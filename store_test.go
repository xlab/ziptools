package ziptools

import (
	"log"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dbPath = "zipcodes.db"
const max = 99999
const min = 10000

func TestGetCity(t *testing.T) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	got, err := db.GetCity(NewZip("13252"))
	assert.NoError(t, err)
	assert.Equal(t, "Syracuse", got)
}

func TestGetLocation(t *testing.T) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	exp := &Location{
		Name:   "Atlanta",
		State:  "TX",
		Locode: NewLocode("TAT"),
	}
	got, err := db.GetLocation(NewLocode("TAT"))
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
}

func TestGetZips(t *testing.T) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	exp := ZipList{
		NewZip("11509"), NewZip("28512"), NewZip("32233"),
	}
	got, err := db.GetZips("Atlantic Beach")
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
}

func TestGetLocodes(t *testing.T) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	exp := LocodeList{
		NewLocode("A2R"), NewLocode("AJI"), NewLocode("ATM"), NewLocode("ATS"),
	}
	got, err := db.GetLocodes("Artesia")
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
}

func TestFindZips(t *testing.T) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	exp := ZipList{
		NewZip("01337"), NewZip("61337"), NewZip("91337"),
	}
	got, err := db.FindZips("1337")
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
}

func TestFindCities(t *testing.T) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	exp := CityList{"Queen Anne", "Princess Anne", "Annemanie", "Saint Anne"}
	got, err := db.FindCities("anne")
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
}

func TestFindLocodes(t *testing.T) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	exp := LocodeList{
		NewLocode("IEB"), NewLocode("IEN"), NewLocode("JBN"),
		NewLocode("JCJ"), NewLocode("LAC"), NewLocode("LC5"),
		NewLocode("LCS"), NewLocode("LNC"), NewLocode("LNS"),
		NewLocode("LNW"), NewLocode("LTX"), NewLocode("LZC"),
		NewLocode("WJF"), NewLocode("ZLI"),
	}
	got, err := db.FindLocodes("lanca")
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
}

// Benchmarks ===============================================

func BenchmarkGetCity(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	zips := make([]Zip, max)
	for i := 0; i < max; i++ {
		zips[i] = NewZip(strconv.Itoa(rand.Intn(max-min) + min))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.GetCity(zips[i%max])
	}
}

func BenchmarkGetLocation(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	locodes := []Locode{
		NewLocode("AON"), NewLocode("ZAM"), NewLocode("ARE"),
		NewLocode("TAT"), NewLocode("AAP"), NewLocode("B8G"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.GetLocation(locodes[i%6])
	}
}

func BenchmarkGetZips(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	cities := []string{
		"New York", "Springfield", "Jamaica",
		"Nesconset", "Richardson", "Netkonfet",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.GetZips(cities[i%6])
	}
}

func BenchmarkGetLocodes(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	cities := []string{
		"New York", "Springfield", "Jamaica",
		"Nesconset", "Richardson", "Netkonfet",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.GetLocodes(cities[i%6])
	}
}

func BenchmarkFindZips(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
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
	defer db.Close()
	cities := []string{
		"york", "field", "spring", "son", "nch", "as",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.FindCities(cities[i%6])
	}
}

func BenchmarkFindLocodes(b *testing.B) {
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	cities := []string{
		"york", "field", "spring", "son", "nch", "as",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.FindLocodes(cities[i%6])
	}
}
