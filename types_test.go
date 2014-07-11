package ziptools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
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
