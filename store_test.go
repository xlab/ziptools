package ziptools

import (
	"log"
	"math/rand"
	"strconv"
	"testing"
)

const dbPath = "zipcodes.db"
const max = 99999
const min = 10000

func BenchmarkGetCity(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	zips := make([]Zip, max)
	for i := 0; i < max; i++ {
		zips[i] = NewZip(strconv.Itoa(rand.Intn(max-min) + min))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.GetCity(zips[i%max])
	}
}

func BenchmarkGetZips(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	cities := []string{
		"New York", "Springfield", "Jamaica",
		"Nesconset", "Richardson", "Netkonfet",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.GetZips(cities[i%6])
	}
}

func BenchmarkFindZips(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	zips := make([]string, max)
	for i := 0; i < max; i++ {
		zips[i] = strconv.Itoa(rand.Intn(min))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.FindZips(zips[i%max])
	}
}

func BenchmarkFindCities(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	cities := []string{
		"york", "field", "spring", "son", "nch", "as",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.FindCities(cities[i%6])
	}
}
