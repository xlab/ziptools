package ziptools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewZip(t *testing.T) {
	zip := NewZip("12345")
	zip2 := NewZip("123")
	assert.Equal(t, "12345", zip.String())
	assert.Equal(t, "123", zip2.String())
}

func TestZipListMarshalJSON(t *testing.T) {
	list := ZipList{
		NewZip("11111"), NewZip("02222"), NewZip("3"),
	}
	exp := `["11111","02222","3"]`
	out, err := json.Marshal(list)
	assert.NoError(t, err)
	assert.Equal(t, exp, string(out))
}

func TestZipListBytes(t *testing.T) {
	list := ZipList{
		NewZip("11111"), NewZip("22222"), NewZip("3"),
	}
	exp := []byte("\x0311111222223\x00\x00\x00\x00")
	assert.Equal(t, exp, list.Bytes())
}

func TestZipListFromBytes(t *testing.T) {
	data := []byte("\x0311111222223\x00\x00\x00\x00")
	exp := ZipList{
		NewZip("11111"), NewZip("22222"), NewZip("3"),
	}
	var list ZipList
	assert.Equal(t, exp, list.FromBytes(data))
}

func TestZipListRange(t *testing.T) {
	data := ZipList{
		NewZip("1"), NewZip("2"), NewZip("3"),
	}
	exp := ZipList{NewZip("2"), NewZip("3")}
	expLimit := ZipList{NewZip("3")}
	assert.Equal(t, exp, data.Range(1, 2))
	assert.Equal(t, expLimit, data.Range(2, 10))
	assert.Empty(t, data.Range(3, 1))
	assert.Empty(t, data.Range(-1, 1))
}

// ==================

func TestNewLocode(t *testing.T) {
	lo := NewLocode("ABD")
	lo2 := NewLocode("AB")
	assert.Equal(t, "ABD", lo.String())
	assert.Equal(t, "AB", lo2.String())
}

func TestLocodeListMarshalJSON(t *testing.T) {
	list := LocodeList{
		NewLocode("ABD"), NewLocode("ABJ"), NewLocode("AQ2"),
	}
	exp := `["ABD","ABJ","AQ2"]`
	out, err := json.Marshal(list)
	assert.NoError(t, err)
	assert.Equal(t, exp, string(out))
}

func TestLocodeListBytes(t *testing.T) {
	list := LocodeList{
		NewLocode("ABD"), NewLocode("ABJ"), NewLocode("AQ"),
	}
	exp := []byte("\x03ABDABJAQ\x00")
	assert.Equal(t, exp, list.Bytes())
}

func TestLocodeListFromBytes(t *testing.T) {
	data := []byte("\x03ABDABJAQ\x00")
	exp := LocodeList{
		NewLocode("ABD"), NewLocode("ABJ"), NewLocode("AQ"),
	}
	var list LocodeList
	assert.Equal(t, exp, list.FromBytes(data))
}

func TestLocodeListRange(t *testing.T) {
	data := LocodeList{
		NewLocode("ABD"), NewLocode("ABJ"), NewLocode("AQ2"),
	}
	exp := LocodeList{NewLocode("ABJ"), NewLocode("AQ2")}
	expLimit := LocodeList{NewLocode("AQ2")}
	assert.Equal(t, exp, data.Range(1, 2))
	assert.Equal(t, expLimit, data.Range(2, 10))
	assert.Empty(t, data.Range(3, 1))
	assert.Empty(t, data.Range(-1, 1))
}

// ==================

func TestLocationBytes(t *testing.T) {
	data := &Location{
		Name:   "Abbeville",
		State:  "AL",
		Locode: NewLocode("ABB"),
	}
	exp := []byte("{\"Name\":\"Abbeville\",\"State\":\"AL\",\"Locode\":\"ABB\"}")
	assert.Equal(t, exp, data.Bytes())
}

func TestLocationFromBytes(t *testing.T) {
	data := []byte("{\"Name\":\"Abbeville\",\"State\":\"AL\",\"Locode\":\"ABB\"}")
	exp := &Location{
		Name:   "Abbeville",
		State:  "AL",
		Locode: NewLocode("ABB"),
	}
	var loc Location
	assert.Equal(t, exp, loc.FromBytes(data))
}
