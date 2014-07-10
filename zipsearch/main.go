// zipsearch tool leverages the ziptools package and provides a simple cli interface for searching zip codes and cities.
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
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/xlab/ziptools"
)

var dbPath string
var cityName bool
var exactMatch bool

func init() {
	flag.BoolVar(&exactMatch, "exact", false, "look for exact match")
	flag.BoolVar(&cityName, "city", false, "given string is a city name or its part")
	flag.StringVar(&dbPath, "db", "zipcodes.db", "specify zip codes database.")
	flag.Parse()
}

func main() {
	if len(flag.Args()) < 1 {
		flag.Usage()
		return
	}
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	db, err := ziptools.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	switch {
	case exactMatch && cityName:
		name := strings.Join(flag.Args(), " ")
		list, err := db.GetZips(name)
		if len(list) < 1 || err != nil {
			fmt.Printf("No zip codes found for %s.\n", name)
			return err
		}
		fmt.Printf("Zip codes in %s: %v\n", name, list)
	case cityName:
		name := strings.Join(flag.Args(), " ")
		list, err := db.FindCities(name)
		if len(list) < 1 || err != nil {
			fmt.Printf("No cities matched %s.\n", name)
			return err
		}
		fmt.Printf("Cities that match %s: %#v\n", name, list)
	case exactMatch:
		zip := ziptools.NewZip(flag.Arg(0))
		city, err := db.GetCity(zip)
		if len(city) < 1 || err != nil {
			fmt.Printf("No city found for %s.\n", zip)
			return err
		}
		fmt.Printf("Zip %s belongs to %s.\n", zip, city)
	default:
		part := flag.Arg(0)
		list, err := db.FindZips(part)
		if len(list) < 1 || err != nil {
			fmt.Printf("No zips matched %s.\n", part)
			return err
		}
		fmt.Printf("Zip codes that match %s: %v\n", part, list)
	}
	return nil
}
