package version

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"

	"github.com/koykov/byteconv"
)

type Compact64 uint64

// ParseCompact64 makes new version from source.
func ParseCompact64(ver []byte) (Compact64, error) {
	return ParseCompact64String(byteconv.B2S(ver))
}

// ParseCompact64String makes new version from source string.
func ParseCompact64String(ver string) (v Compact64, err error) {
	err = v.ParseString(ver)
	return
}

// NewCompact64 composes version from given parts.
func NewCompact64(major, minor, patch, revision uint16) Compact64 {
	var v Compact64
	v = v | Compact64(major)<<48
	v = v | Compact64(minor)<<32
	v = v | Compact64(patch)<<16
	v = v | Compact64(revision)
	return v
}

func (v *Compact64) Parse(ver []byte) error {
	return v.ParseString(byteconv.B2S(ver))
}

func (v *Compact64) ParseString(ver string) error {
	if len(ver) == 0 {
		return ErrEmpty
	}

	var m, n, p, r uint16
	c := 0
	for {
		i := strings.Index(ver, ".")
		if i == -1 {
			break
		}
		raw := ver[:i]
		u, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return err
		}
		switch c {
		case 0:
			m = uint16(u)
		case 1:
			n = uint16(u)
		case 2:
			p = uint16(u)
		}
		c++
		ver = ver[i+1:]
	}
	u, err := strconv.ParseUint(ver, 10, 64)
	if err != nil {
		return err
	}
	switch c {
	case 0:
		m = uint16(u)
	case 1:
		n = uint16(u)
	case 2:
		p = uint16(u)
	case 3:
		r = uint16(u)
	}
	v.SetMajor(m).SetMinor(n).SetPatch(p).SetRevision(r)
	return nil
}

func (v *Compact64) SetMajor(value uint16) *Compact64 {
	*v = *v | Compact64(value)<<48
	return v
}

func (v *Compact64) SetMinor(value uint16) *Compact64 {
	*v = *v | Compact64(value)<<32
	return v
}

func (v *Compact64) SetPatch(value uint16) *Compact64 {
	*v = *v | Compact64(value)<<16
	return v
}

func (v *Compact64) SetRevision(value uint16) *Compact64 {
	*v = *v | Compact64(value)
	return v
}

func (v *Compact64) Major() uint16 {
	return uint16(*v >> 48)
}

func (v *Compact64) Minor() uint16 {
	return uint16(*v >> 32)
}

func (v *Compact64) Patch() uint16 {
	return uint16(*v >> 16)
}

func (v *Compact64) Revision() uint16 {
	return uint16(*v)
}

func (v *Compact64) Bytes() []byte {
	return byteconv.S2B(v.String())
}

func (v *Compact64) String() string {
	var a [23]byte
	buf := a[:][:0]
	m, n, p, r := v.Major(), v.Minor(), v.Patch(), v.Revision()
	switch {
	case r > 0:
		buf = strconv.AppendUint(buf, uint64(m), 10)
		buf = append(buf, '.')
		buf = strconv.AppendUint(buf, uint64(n), 10)
		buf = append(buf, '.')
		buf = strconv.AppendUint(buf, uint64(p), 10)
		buf = append(buf, '.')
		buf = strconv.AppendUint(buf, uint64(r), 10)
	case r == 0 && p > 0:
		buf = strconv.AppendUint(buf[:0], uint64(m), 10)
		buf = append(buf, '.')
		buf = strconv.AppendUint(buf, uint64(n), 10)
		buf = append(buf, '.')
		buf = strconv.AppendUint(buf, uint64(p), 10)
	case r == 0 && p == 0:
		buf = strconv.AppendUint(buf[:0], uint64(m), 10)
		buf = append(buf, '.')
		buf = strconv.AppendUint(buf, uint64(n), 10)
	}
	return byteconv.B2S(buf)
}

func (v *Compact64) WriteBinaryTo(w io.Writer) (int64, error) {
	p, _ := v.MarshalBinary()
	n, err := w.Write(p)
	return int64(n), err
}

func (v *Compact64) MarshalBinary() ([]byte, error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(*v))
	return buf[:], nil
}

func (v *Compact64) MarshalText() ([]byte, error) {
	return v.Bytes(), nil
}

func (v *Compact64) UnmarshalBinary(bin []byte) error {
	if len(bin) < 8 {
		return ErrBinLen8
	}
	*v = Compact64(binary.LittleEndian.Uint64(bin))
	return nil
}

func (v *Compact64) UnmarshalText(p []byte) error {
	return v.Parse(p)
}

var _, _ = ParseCompact64, NewCompact64
