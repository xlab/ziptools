package ziptools

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"
)

const (
	// ZipLen is the default length of Zip
	ZipLen = 5
	// LocodeLen is the default length of Locode
	LocodeLen = 3
)

// Zip represents a zip code.
type Zip [ZipLen]byte

// Locode represents a locode.
type Locode [LocodeLen]byte

// Location represents transport location.
type Location struct {
	Name   string
	State  string
	Locode Locode
}

// Bytes returns a serialized version of a location.
func (l Location) Bytes() []byte {
	b, _ := json.Marshal(l)
	return b
}

// FromBytes constructs a new location from bytes.
func (l *Location) FromBytes(b []byte) *Location {
	json.Unmarshal(b, l)
	return l
}

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

// NewLocode creates a new locode from string.
func NewLocode(str string) (code Locode) {
	for i, c := range []byte(str) {
		if i >= len(code) {
			return
		}
		code[i] = c
	}
	return
}

// String represents a zip as a string.
func (z Zip) String() string {
	return string(z.Bytes())
}

// String represents a locode as a string.
func (l Locode) String() string {
	return string(l.Bytes())
}

// MarshalJSON represents a zip as a string while marshaling as JSON.
func (z Zip) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(z.String())), nil
}

// UnmarshalJSON restores zip code from bytes after marshaling as JSON.
func (z *Zip) UnmarshalJSON(b []byte) (err error) {
	str := string(b)
	if str, err = strconv.Unquote(str); err != nil {
		return err
	}
	*z = NewZip(str)
	return
}

// MarshalJSON represents a locode as a string while marshaling as JSON.
func (l Locode) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(l.String())), nil
}

// UnmarshalJSON restores locode from bytes after marshaling as JSON.
func (l *Locode) UnmarshalJSON(b []byte) (err error) {
	str := string(b)
	if str, err = strconv.Unquote(str); err != nil {
		return err
	}
	*l = NewLocode(str)
	return
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

// Bytes represents a locode as bytes.
func (l Locode) Bytes() []byte {
	b := make([]byte, 0, LocodeLen)
	for i := range l {
		if l[i] == 0 {
			return b
		}
		b = append(b, l[i])
	}
	return b
}

// ZipList reperesents a list of zip codes.
type ZipList []Zip

// LocodeList reperesents a list of locodes.
type LocodeList []Locode

// CityList reperesents a list of cities.
type CityList []string

// Range returns a sliced variant of a zip list.
func (z ZipList) Range(offset, limit int) ZipList {
	if offset < 0 || offset >= len(z) {
		return ZipList{}
	}
	if offset+limit > len(z) {
		limit = len(z) - offset
	}
	return z[offset : offset+limit]
}

// Range returns a sliced variant of a locode list.
func (l LocodeList) Range(offset, limit int) LocodeList {
	if offset < 0 || offset >= len(l) {
		return LocodeList{}
	}
	if offset+limit > len(l) {
		limit = len(l) - offset
	}
	return l[offset : offset+limit]
}

// Range returns a sliced variant of a city list.
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

// Bytes returns a serialized version of a locode list. First byte represents the length.
//
//  [N][locode1][locode2]...[locodeN]
func (l LocodeList) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(len(l)))
	binary.Write(&buf, binary.LittleEndian, l)
	return buf.Bytes()
}

// FromBytes constructs a new zip list from bytes. First byte reperesents the length.
func (z *ZipList) FromBytes(b []byte) ZipList {
	r := bytes.NewReader(b)
	if n, err := r.ReadByte(); err != nil {
		return nil
	} else {
		*z = make(ZipList, n)
	}
	binary.Read(r, binary.LittleEndian, z)
	return *z
}

// FromBytes constructs a new locode from bytes. First byte reperesents the length.
func (l *LocodeList) FromBytes(b []byte) LocodeList {
	r := bytes.NewReader(b)
	if n, err := r.ReadByte(); err != nil {
		return nil
	} else {
		*l = make(LocodeList, n)
	}
	binary.Read(r, binary.LittleEndian, l)
	return *l
}
