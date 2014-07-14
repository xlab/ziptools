package ziptools

import (
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	citiesBuck     = []byte("cities")
	locodesBuck    = []byte("locodes")
	locationsBuck  = []byte("locations")
	zipsBuck       = []byte("zips")
	subZipsBuck    = []byte("subzips")
	subCitiesBuck  = []byte("subcities")
	subLocodesBuck = []byte("sublocodes")
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

// GetCity gets a city that is assigned to the specified zip code.
// This methods looks for an exact match.
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

// GetLocation gets a location that is assigned to the specified locode.
// This methods looks for an exact match.
func (d *DB) GetLocation(l Locode) (loc *Location, err error) {
	loc = &Location{}
	err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(locationsBuck); b != nil {
			loc.FromBytes(b.Get(l.Bytes()))
			return nil
		}
		return bolt.ErrBucketNotFound
	})
	return
}

// Get a list of zip codes in the specified city. This methods looks
// for an exact match.
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

// Get a list of locodes for the specified city. This methods looks
// for an exact match.
func (d *DB) GetLocodes(city string) (locodes LocodeList, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(locodesBuck); b != nil {
			locodes.FromBytes(b.Get([]byte(city)))
			return nil
		}
		return bolt.ErrBucketNotFound
	})
	return
}

// Find all cities that match the given substring.
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

// Find all locodes by a given substring of a city name.
func (d *DB) FindLocodes(citypart string) (locodes LocodeList, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(subLocodesBuck); b != nil {
			locodes.FromBytes(b.Get([]byte(citypart)))
			return nil
		}
		return bolt.ErrBucketNotFound
	})
	return
}

// Find all zip codes that match the given substring.
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
