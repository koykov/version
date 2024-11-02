package version

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"

	"github.com/koykov/byteconv"
)

// Compact32 represents simple version based on uint32 type.
type Compact32 uint32

// ParseCompact32 makes new version from source.
func ParseCompact32(ver []byte) (Compact32, error) {
	return ParseCompact32String(byteconv.B2S(ver))
}

// ParseCompact32String makes new version from source string.
func ParseCompact32String(ver string) (v Compact32, err error) {
	err = v.ParseString(ver)
	return
}

// NewCompact32 composes version from given parts.
func NewCompact32(major, minor, patch, revision uint8) Compact32 {
	var v Compact32
	v = v | Compact32(major)<<24
	v = v | Compact32(minor)<<16
	v = v | Compact32(patch)<<8
	v = v | Compact32(revision)
	return v
}

func (v *Compact32) Parse(ver []byte) error {
	return v.ParseString(byteconv.B2S(ver))
}

func (v *Compact32) ParseString(ver string) error {
	if len(ver) == 0 {
		return ErrEmpty
	}

	var m, n, p, r uint8
	c := 0
	for {
		i := strings.Index(ver, ".")
		if i == -1 {
			break
		}
		raw := ver[:i]
		u, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return err
		}
		switch c {
		case 0:
			m = uint8(u)
		case 1:
			n = uint8(u)
		case 2:
			p = uint8(u)
		}
		c++
		ver = ver[i+1:]
	}
	u, err := strconv.ParseUint(ver, 10, 32)
	if err != nil {
		return err
	}
	switch c {
	case 0:
		m = uint8(u)
	case 1:
		n = uint8(u)
	case 2:
		p = uint8(u)
	case 3:
		r = uint8(u)
	}
	v.SetMajor(m).SetMinor(n).SetPatch(p).SetRevision(r)
	return nil
}

func (v *Compact32) SetMajor(value uint8) *Compact32 {
	*v = *v | Compact32(value)<<24
	return v
}

func (v *Compact32) SetMinor(value uint8) *Compact32 {
	*v = *v | Compact32(value)<<16
	return v
}

func (v *Compact32) SetPatch(value uint8) *Compact32 {
	*v = *v | Compact32(value)<<8
	return v
}

func (v *Compact32) SetRevision(value uint8) *Compact32 {
	*v = *v | Compact32(value)
	return v
}

func (v *Compact32) Major() uint8 {
	return uint8(*v >> 24)
}

func (v *Compact32) Minor() uint8 {
	return uint8(*v >> 16)
}

func (v *Compact32) Patch() uint8 {
	return uint8(*v >> 8)
}

func (v *Compact32) Revision() uint8 {
	return uint8(*v)
}

func (v *Compact32) AppendBytes(dst []byte) []byte {
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

func (v *Compact32) Bytes() (r []byte) {
	var buf [24]byte
	r = v.AppendBytes(buf[:])
	return
}

func (v *Compact32) String() string {
	return byteconv.B2S(v.Bytes())
}

func (v *Compact32) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(v.Bytes())
	return int64(n), err
}

func (v *Compact32) WriteBinaryTo(w io.Writer) (int32, error) {
	p, _ := v.MarshalBinary()
	n, err := w.Write(p)
	return int32(n), err
}

func (v *Compact32) MarshalBinary() ([]byte, error) {
	var buf [8]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(*v))
	return buf[:], nil
}

func (v *Compact32) MarshalText() ([]byte, error) {
	return v.Bytes(), nil
}

func (v *Compact32) UnmarshalBinary(bin []byte) error {
	if len(bin) < 8 {
		return ErrBinLen8
	}
	*v = Compact32(binary.LittleEndian.Uint32(bin))
	return nil
}

func (v *Compact32) UnmarshalText(p []byte) error {
	return v.Parse(p)
}

var _, _ = ParseCompact32, NewCompact32