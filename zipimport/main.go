// zipimport tool is suited for Bolt DB creation from a gzipped CSV with zip codes.
// This operation may take a few minutes.
//
//   $ zipimport -h
//   Usage of zipimport:
//     -csv="zip_code_database.csv.gz": gzipped .csv file with zip codes.
//     -db="zipcodes.db": file to store a newly created zip codes database.
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/xlab/ziptools"
)

var (
	citiesBuck    = []byte("cities")
	zipsBuck      = []byte("zips")
	subZipsBuck   = []byte("subzips")
	subCitiesBuck = []byte("subcities")
)

var dbPath string
var csvPath string

func init() {
	flag.StringVar(&dbPath, "db", "zipcodes.db", "file to store a newly created zip codes database.")
	flag.StringVar(&csvPath, "csv", "zip_code_database.csv.gz", "gzipped .csv file with zip codes.")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

type DB struct {
	db *bolt.DB
}

func run() (err error) {
	db := new(DB)

	// open the DB file
	if db.db, err = bolt.Open(dbPath, 0644, nil); err != nil {
		return
	}
	defer db.db.Close()

	// open the CSV file
	gzipf, err := os.Open(csvPath)
	if err != nil {
		return
	}
	defer gzipf.Close()
	// decompress stream
	var r io.Reader
	if r, err = gzip.NewReader(gzipf); err != nil {
		return err
	}
	var n int
	if n, err = db.addZips(csv.NewReader(r)); err != nil {
		return
	}
	log.Printf("zipimport: %d zip codes imported", n)
	if err = db.addSubstrings(); err != nil {
		return
	}
	log.Println("zipimport: done indexing")
	return
}

func (d *DB) addZips(csv *csv.Reader) (n int, err error) {
	// begin a writing transaction
	tx, err := d.db.Begin(true)
	if err != nil {
		return
	}
	var zips *bolt.Bucket
	if zips, err = tx.CreateBucketIfNotExists(zipsBuck); err != nil {
		return
	}
	// read CSV and fill the buckets
	for {
		var fields []string
		fields, err = csv.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("zipimport: ignored a csv line due to an error", err)
		}
		if fields[1] == "MILITARY" {
			continue
		}
		zip := ziptools.NewZip(fields[0])
		// zip = city
		if err = zips.Put(zip.Bytes(), []byte(fields[2])); err != nil {
			return
		}
		n++
	}
	return n, tx.Commit()
}

func (d *DB) addSubstrings() (err error) {
	// begin a writing transaction
	tx, err := d.db.Begin(true)
	if err != nil {
		return
	}
	// create buckets
	var cities *bolt.Bucket
	var subcities *bolt.Bucket
	var subzips *bolt.Bucket
	if cities, err = tx.CreateBucketIfNotExists(citiesBuck); err != nil {
		return
	}
	if subcities, err = tx.CreateBucketIfNotExists(subCitiesBuck); err != nil {
		return
	}
	if subzips, err = tx.CreateBucketIfNotExists(subZipsBuck); err != nil {
		return
	}

	errC := make(chan error, 1)
	pairs := make(chan struct{ k, v []byte }, 100)
	go func() {
		seen := make(map[string]struct{})
		// this is a writing goroutine
		for p := range pairs {
			zip := string(p.k)
			city := string(bytes.ToLower(p.v))
			// put full city name -> ziplist
			list := d.getList(cities, p.v)
			list = append(list, ziptools.NewZip(zip))
			if err := cities.Put(p.v, list.Bytes()); err != nil {
				errC <- err
				return
			}
			// put subzips -> ziplist
			if err = d.putSubstringZipList(subzips, zip, ziptools.NewZip(zip)); err != nil {
				errC <- err
				return
			}
			// put subcities -> ziplist
			// cities are not unique, so filter
			if _, ok := seen[city]; ok {
				continue
			}
			seen[city] = struct{}{}
			if err = d.putSubstringZipList(subcities, city, ziptools.NewZip(zip)); err != nil {
				errC <- err
				return
			}
		}
		errC <- nil
	}()

	// Iterate over zip codes in read-only tx
	if err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(zipsBuck); b != nil {
			err := b.ForEach(func(k []byte, v []byte) error {
				select {
				case err := <-errC:
					return err
				default:
					pairs <- struct{ k, v []byte }{k, v}
					return nil
				}
			})
			close(pairs)
			return err
		}
		return bolt.ErrBucketNotFound
	}); err != nil {
		return
	}

	// wait until writer is done
	<-errC
	return tx.Commit()
}

// putSubstringZipList generates all possible substrings (prepend, append),
// and puts them to bucket as keys to ZipLists.
func (d *DB) putSubstringZipList(buck *bolt.Bucket, str string, zip ziptools.Zip) error {
	seen := make(map[string]struct{})
	put := func(substr string) error {
		if _, ok := seen[substr]; ok || len(substr) < 1 {
			return nil
		}
		seen[substr] = struct{}{}
		list := d.getList(buck, []byte(substr))
		list = append(list, zip)
		if err := buck.Put([]byte(substr), list.Bytes()); err != nil {
			return err
		}
		return nil
	}

	for i := range str {
		if err := put(str[0:i]); err != nil {
			return err
		}
	}
	for i := range str {
		if err := put(str[i:len(str)]); err != nil {
			return err
		}
	}
	return nil
}

// Gets a ZipList by key from a bucket.
func (d *DB) getList(buck *bolt.Bucket, key []byte) (list ziptools.ZipList) {
	return list.FromBytes(buck.Get(key))
}
