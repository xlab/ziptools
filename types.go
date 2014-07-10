// Package ziptools provides functionality to search through USA zip codes and cities extremely fast.
// Internally it uses the Bolt — blazing fast in-memory key-value storage that allows to
// cache the codes and city names effectively.
//
// 	$ go test -bench=.
// 	PASS
// 	BenchmarkGetCity	 1000000	      2587 ns/op
// 	BenchmarkGetZips	  200000	     12890 ns/op
// 	BenchmarkFindZips	  200000	      7123 ns/op
// 	BenchmarkFindCities	    5000	    358353 ns/op
// 	ok  	github.com/xlab/ziptools	9.562s
//
// Database should be created using a CSV file located at http://www.unitedstateszipcodes.org/zip_code_database.csv.
// The gzipped version of that file with stripped CSV header is included within this package.
//
// The zipimport tool is suited for Bolt DB creation from that gzipped CSV.
//
//   $ zipimport -h
//   Usage of zipimport:
//     -csv="zip_code_database.csv.gz": gzipped .csv file with zip codes.
//     -db="zipcodes.db": file to store a newly created zip codes database.
//
// The zipsearch tool leverages the ziptools package and provides a simple cli interface
// for searching zip codes and cities within a console window (for testing purposes).
//
// Usage of zipsearch:
//   -city=false: given string is a city name or its part
//   -db="zipcodes.db": specify zip codes database.
//   -exact=false: look for exact match
// List all zipcodes in city:
//   $ zipsearch -exact -city Richardson
//   Zip codes in Richardson: [75080 75081 75082 75083 75085]
//
// Get the city that has the specified zip:
//   $ zipsearch -exact 10106
//   Zip 10106 belongs to New York.
//
// List all cities that match the given substring:
//   $ zipsearch -city english
//   Cities that match english: ziptools.CityList{"Englishtown", "English", "North English", "South English"}
//
// List all zips that match the given substring:
//   $ zipsearch 1337
//   Zip codes that match 1337: [01337 61337 91337]
package ziptools

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	// ZipLen is the default length of Zip
	ZipLen = 5
)

// Zip represents a zip code.
type Zip [ZipLen]byte

// NewZip creates a new zip code from string.
func NewZip(str string) (zip Zip) {
	for i, c := range []byte(str) {
		if i >= len(zip) {
			return
		}
		zip[i] = c
	}
	return
}

// String represents a zip as a string.
func (z Zip) String() string {
	return fmt.Sprintf("%s", z.Bytes())
}

// Bytes represents a zip as bytes.
func (z Zip) Bytes() []byte {
	b := make([]byte, 0, ZipLen)
	for i := range z {
		if z[i] == 0 {
			return b
		}
		b = append(b, z[i])
	}
	return b
}

// List reperesents a list of zip codes.
type ZipList []Zip

// List reperesents a list of cities.
type CityList []string

// Range returns a sliced variant of the zip list.
func (z ZipList) Range(offset, limit int) ZipList {
	if offset < 0 || offset >= len(z) {
		return ZipList{}
	}
	if offset+limit > len(z) {
		limit = len(z) - offset
	}
	return z[offset : offset+limit]
}

// Range returns a sliced variant of the city list.
func (c CityList) Range(offset, limit int) CityList {
	if offset < 0 || offset >= len(c) {
		return CityList{}
	}
	if offset+limit > len(c) {
		limit = len(c) - offset
	}
	return c[offset : offset+limit]
}

// Bytes returns a serialized version of a zip list. First byte represents the length.
//
// [N][zip1][zip2]...[zipN]
func (z ZipList) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(len(z)))
	binary.Write(&buf, binary.LittleEndian, z)
	return buf.Bytes()
}

// FromBytes constructs a new zip list. First byte reperesents the length.
func (z *ZipList) FromBytes(b []byte) ZipList {
	r := bytes.NewReader(b)
	if l, err := r.ReadByte(); err != nil {
		return nil
	} else {
		*z = make(ZipList, l)
	}
	binary.Read(r, binary.LittleEndian, z)
	return *z
}
