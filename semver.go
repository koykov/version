package version

import (
	"encoding/binary"
	"errors"
	"io"
	"strconv"

	"github.com/koykov/byteconv"
)

type Semver struct {
	m, n, p  uint32
	pr, meta string
}

// ParseSemver makes new version from source.
func ParseSemver(ver []byte) (Semver, error) {
	return ParseSemverString(byteconv.B2S(ver))
}

// ParseSemverString makes new version from source string.
func ParseSemverString(ver string) (v Semver, err error) {
	err = v.ParseString(ver)
	return
}

func (v *Semver) Parse(ver []byte) error {
	return v.ParseString(byteconv.B2S(ver))
}

func (v *Semver) ParseString(ver string) error {
	const (
		partMajor = iota
		partMinor
		partPatch
		partPreRelease
		partMeta
	)

	n := len(ver)
	if n == 0 {
		return ErrEmptySemver
	}
	_ = ver[n-1]
	var offset, part int

	for ; ver[offset] == ' ' || ver[offset] == 'v'; offset++ {
	}
	for ; n > offset && ver[n-1] == ' '; n-- {
	}

	var i int
	for i = offset; i < n; i++ {
		if ver[i] == '.' {
			x, err := strconv.ParseUint(ver[offset:i], 10, 32)
			if err != nil {
				return err
			}
			switch part {
			case partMajor:
				v.m = uint32(x)
				part = partMinor
			case partMinor:
				v.n = uint32(x)
				part = partPatch
			default:
				return ErrBadSemver
			}
			offset = i + 1
		}
		if ver[i] == '-' && part == partPatch {
			x, err := strconv.ParseUint(ver[offset:i], 10, 32)
			if err != nil {
				return err
			}
			v.p = uint32(x)
			part = partPreRelease
		}
		if ver[i] == '+' && part == partPreRelease {
			v.pr = ver[offset:i]
			part = partMeta
		}
	}

	switch {
	case n-offset > 0 && part < partPreRelease:
		x, err := strconv.ParseUint(ver[offset:i], 10, 32)
		if err != nil {
			return err
		}
		switch part {
		case partMinor:
			v.n = uint32(x)
		case partMajor:
			v.p = uint32(x)
		default:
			return ErrBadSemver
		}
	case n-offset > 0 && part == partPreRelease:
	}
	return nil
}

func (v *Semver) SetMajor(value uint32) *Semver {
	v.m = value
	return v
}

func (v *Semver) SetMinor(value uint32) *Semver {
	v.n = value
	return v
}

func (v *Semver) SetPatch(value uint32) *Semver {
	v.p = value
	return v
}

func (v *Semver) SetPreRelease(value []byte) *Semver {
	v.pr = byteconv.B2S(value)
	return v
}

func (v *Semver) SetPreReleaseString(value string) *Semver {
	v.pr = value
	return v
}

func (v *Semver) SetMeta(value []byte) *Semver {
	v.meta = byteconv.B2S(value)
	return v
}

func (v *Semver) SetMetaString(value string) *Semver {
	v.meta = value
	return v
}

func (v *Semver) Major() uint32 {
	return v.m
}

func (v *Semver) Minor() uint32 {
	return v.n
}

func (v *Semver) Patch() uint32 {
	return v.p
}

func (v *Semver) PreRelease() []byte {
	return byteconv.S2B(v.pr)
}

func (v *Semver) PreReleaseString() string {
	return v.pr
}

func (v *Semver) Meta() []byte {
	return byteconv.S2B(v.meta)
}

func (v *Semver) MetaString() string {
	return v.meta
}

func (v *Semver) MarshalBinary() (data []byte, err error) {
	var a [128]byte
	buf := a[:]
	binary.LittleEndian.PutUint32(buf[:4], v.m)
	binary.LittleEndian.PutUint32(buf[4:], v.n)
	binary.LittleEndian.PutUint32(buf[8:], v.p)
	binary.LittleEndian.PutUint32(buf[12:], uint32(len(v.pr)))
	buf = append(buf, v.pr...)
	binary.LittleEndian.PutUint32(buf[len(buf)-1:], uint32(len(v.meta)))
	buf = append(buf, v.meta...)
	return buf, nil
}

func (v *Semver) UnmarshalBinary(data []byte) error {
	// todo implement me
	return nil
}

func (v *Semver) AppendBytes(dst []byte) []byte {
	// todo implement me
	return dst
}

func (v *Semver) MarshalText() (text []byte, err error) {
	return v.Bytes(), nil
}

func (v *Semver) UnmarshalText(text []byte) error {
	return v.Parse(text)
}

func (v *Semver) Bytes() []byte {
	var buf [128]byte
	return v.AppendBytes(buf[:0])
}

func (v *Semver) String() string {
	return byteconv.B2S(v.Bytes())
}

func (v *Semver) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(v.Bytes())
	return int64(n), err
}

func (v *Semver) WriteBinaryTo(w io.Writer) (int64, error) {
	p, _ := v.MarshalBinary()
	n, err := w.Write(p)
	return int64(n), err
}

var (
	ErrEmptySemver = errors.New("version is empty")
	ErrBadSemver   = errors.New("wrong semver format")
)
