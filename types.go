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

// MarshalJSON represents a zip as a string while marshaling as JSON.
func (z Zip) MarshalJSON() ([]byte, error) {
	return []byte(z.String()), nil
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
//  [N][zip1][zip2]...[zipN]
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
