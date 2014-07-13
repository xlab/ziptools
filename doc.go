// Package ziptools provides functionality to search through USA zip codes and cities in extremely fast way.
// Internally it uses the Bolt — blazing fast in-memory key-value storage that allows to
// cache the codes and city names effectively. LOCODE searches supported.
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
//     -zips="zip_code_database.csv.gz": gzipped .csv file with zip codes.
//     -locodes="us_locode_database.csv.gz": gzipped .csv file with locodes.
//     -db="zipcodes.db": file to store a newly created zip codes database.
//
// Installation and Examples
//
// After the Bolt database is created, you may remove zip_code_database.csv.gz.
//
//   go get https://github.com/xlab/ziptools/zipimport
//   go get https://github.com/xlab/ziptools/zipsearch
//   zipimport
//   du -csh zipcodes.db
//
// The zipsearch tool leverages the ziptools package and provides a simple cli interface
// for searching zip codes and cities within a console window (for testing purposes).
//
//   $ zipsearch -h
//   Usage of zipsearch:
//     -city=false: given string is a city name or its part
//     -db="zipcodes.db": specify zip codes database.
//     -exact=false: look for exact match
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
