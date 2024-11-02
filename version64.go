package version

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"

	"github.com/koykov/byteconv"
)

// Version64 represents simple version based on uint64 type.
type Version64 uint64

// ParseVersion64 makes new version from source.
func ParseVersion64(ver []byte) (Version64, error) {
	return ParseVersion64String(byteconv.B2S(ver))
}

// ParseVersion64String makes new version from source string.
func ParseVersion64String(ver string) (v Version64, err error) {
	err = v.ParseString(ver)
	return
}

// NewCompact64 composes version from given parts.
func NewCompact64(major, minor, patch, revision uint16) Version64 {
	var v Version64
	v = v | Version64(major)<<48
	v = v | Version64(minor)<<32
	v = v | Version64(patch)<<16
	v = v | Version64(revision)
	return v
}

func (v *Version64) Parse(ver []byte) error {
	return v.ParseString(byteconv.B2S(ver))
}

func (v *Version64) ParseString(ver string) error {
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

func (v *Version64) SetMajor(value uint16) *Version64 {
	*v = *v | Version64(value)<<48
	return v
}

func (v *Version64) SetMinor(value uint16) *Version64 {
	*v = *v | Version64(value)<<32
	return v
}

func (v *Version64) SetPatch(value uint16) *Version64 {
	*v = *v | Version64(value)<<16
	return v
}

func (v *Version64) SetRevision(value uint16) *Version64 {
	*v = *v | Version64(value)
	return v
}

func (v *Version64) Major() uint16 {
	return uint16(*v >> 48)
}

func (v *Version64) Minor() uint16 {
	return uint16(*v >> 32)
}

func (v *Version64) Patch() uint16 {
	return uint16(*v >> 16)
}

func (v *Version64) Revision() uint16 {
	return uint16(*v)
}

func (v *Version64) AppendBytes(dst []byte) []byte {
	m, n, p, r := v.Major(), v.Minor(), v.Patch(), v.Revision()
	switch {
	case r > 0:
		dst = strconv.AppendUint(dst, uint64(m), 10)
		dst = append(dst, '.')
		dst = strconv.AppendUint(dst, uint64(n), 10)
		dst = append(dst, '.')
		dst = strconv.AppendUint(dst, uint64(p), 10)
		dst = append(dst, '.')
		dst = strconv.AppendUint(dst, uint64(r), 10)
	case r == 0 && p > 0:
		dst = strconv.AppendUint(dst[:0], uint64(m), 10)
		dst = append(dst, '.')
		dst = strconv.AppendUint(dst, uint64(n), 10)
		dst = append(dst, '.')
		dst = strconv.AppendUint(dst, uint64(p), 10)
	case r == 0 && p == 0:
		dst = strconv.AppendUint(dst[:0], uint64(m), 10)
		dst = append(dst, '.')
		dst = strconv.AppendUint(dst, uint64(n), 10)
	}
	return dst
}

func (v *Version64) Bytes() (r []byte) {
	var buf [24]byte
	r = v.AppendBytes(buf[:0])
	return
}

func (v *Version64) String() string {
	return byteconv.B2S(v.Bytes())
}

func (v *Version64) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(v.Bytes())
	return int64(n), err
}

func (v *Version64) WriteBinaryTo(w io.Writer) (int64, error) {
	p, _ := v.MarshalBinary()
	n, err := w.Write(p)
	return int64(n), err
}

func (v *Version64) MarshalBinary() ([]byte, error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(*v))
	return buf[:], nil
}

func (v *Version64) MarshalText() ([]byte, error) {
	return v.Bytes(), nil
}

func (v *Version64) UnmarshalBinary(bin []byte) error {
	if len(bin) < 8 {
		return ErrBinLen8
	}
	*v = Version64(binary.LittleEndian.Uint64(bin))
	return nil
}

func (v *Version64) UnmarshalText(p []byte) error {
	return v.Parse(p)
}

var _, _ = ParseVersion64, NewCompact64
