package ziptools

import (
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	citiesBuck    = []byte("cities")
	zipsBuck      = []byte("zips")
	subZipsBuck   = []byte("subzips")
	subCitiesBuck = []byte("subcities")
)

// DB abstracts database access.
type DB struct {
	db *bolt.DB
}

// Open opens a Bolt database from a file if it exists.
func Open(path string) (db *DB, err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}
	db = new(DB)
	db.db, err = bolt.Open(path, 0644, nil)
	return
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) GetCity(z Zip) (city string, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(zipsBuck); b != nil {
			city = string(b.Get(z.Bytes()))
			return nil
		}
		return bolt.ErrBucketNotFound
	})
	return
}

func (d *DB) GetZips(city string) (zips ZipList, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(citiesBuck); b != nil {
			zips.FromBytes(b.Get([]byte(city)))
			return nil
		}
		return bolt.ErrBucketNotFound
	})
	return
}

func (d *DB) FindCities(citypart string) (cities CityList, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(subCitiesBuck); b != nil {
			var list ZipList
			citypart = strings.ToLower(citypart)
			list.FromBytes(b.Get([]byte(citypart)))
			for _, zip := range list {
				if city, err := d.GetCity(zip); err != nil {
					return err
				} else {
					cities = append(cities, city)
				}
			}
			return nil
		}
		return bolt.ErrBucketNotFound
	})
	return
}

func (d *DB) FindZips(zippart string) (zips ZipList, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(subZipsBuck); b != nil {
			zips.FromBytes(b.Get([]byte(zippart)))
			return nil
		}
		return bolt.ErrBucketNotFound
	})
	return
}
