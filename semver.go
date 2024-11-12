package version

import (
	"encoding/binary"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/koykov/byteconv"
)

type Semver struct {
	m, n, p   uint64
	pre, meta string
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

func NewSemver(major, minor, patch uint64) *Semver {
	return &Semver{
		m: major,
		n: minor,
		p: patch,
	}
}

func (v *Semver) WithPreRelease(pre string) *Semver {
	v.pre = pre
	return v
}

func (v *Semver) WithMeta(meta string) *Semver {
	v.meta = meta
	return v
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
				v.m = x
				part = partMinor
			case partMinor:
				v.n = x
				part = partPatch
			default:
				return ErrBadSemver
			}
			offset = i + 1
		}
		if ver[i] == '-' {
			x, err := strconv.ParseUint(ver[offset:i], 10, 32)
			if err != nil {
				return err
			}
			v.p = x
			part = partPreRelease
			offset = i + 1
			break
		}
		if ver[i] == '+' {
			x, err := strconv.ParseUint(ver[offset:i], 10, 32)
			if err != nil {
				return err
			}
			v.p = x
			part = partMeta
			offset = i + 1
			break
		}
	}

	switch part {
	case partPatch:
		x, err := strconv.ParseUint(ver[offset:n], 10, 32)
		if err != nil {
			return err
		}
		v.p = x
	case partPreRelease:
		v.pre = ver[offset:]
		if i = strings.IndexByte(ver[offset:], '+'); i != -1 {
			v.pre = ver[offset : offset+i]
			v.meta = ver[offset+i+1:]
			return nil
		}
	case partMeta:
		v.meta = ver[offset:]
	default:
		return ErrBadSemver
	}

	return nil
}

func (v *Semver) SetMajor(value uint64) *Semver {
	v.m = value
	return v
}

func (v *Semver) SetMinor(value uint64) *Semver {
	v.n = value
	return v
}

func (v *Semver) SetPatch(value uint64) *Semver {
	v.p = value
	return v
}

func (v *Semver) SetPreRelease(value []byte) *Semver {
	v.pre = byteconv.B2S(value)
	return v
}

func (v *Semver) SetPreReleaseString(value string) *Semver {
	v.pre = value
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

func (v *Semver) Major() uint64 {
	return v.m
}

func (v *Semver) Minor() uint64 {
	return v.n
}

func (v *Semver) Patch() uint64 {
	return v.p
}

func (v *Semver) PreRelease() []byte {
	return byteconv.S2B(v.pre)
}

func (v *Semver) PreReleaseString() string {
	return v.pre
}

func (v *Semver) Meta() []byte {
	return byteconv.S2B(v.meta)
}

func (v *Semver) MetaString() string {
	return v.meta
}

func (v *Semver) MarshalBinary() (data []byte, err error) {
	buf := make([]byte, len(v.pre)+len(v.meta)+32) // 32 == uint32(len(pre)) + uint32(len(meta)) + v.m + v.n + v.p
	binary.LittleEndian.PutUint64(buf[:8], v.m)
	binary.LittleEndian.PutUint64(buf[8:], v.n)
	binary.LittleEndian.PutUint64(buf[16:], v.p)
	binary.LittleEndian.PutUint32(buf[24:], uint32(len(v.pre)))
	buf = append(buf[:28], v.pre...)
	binary.LittleEndian.PutUint32(buf[len(buf)-1:], uint32(len(v.meta)))
	buf = append(buf, v.meta...)
	return buf, nil
}

func (v *Semver) UnmarshalBinary(data []byte) error {
	// todo implement me
	return nil
}

func (v *Semver) AppendBytes(dst []byte) []byte {
	dst = strconv.AppendUint(dst, v.m, 10)
	dst = append(dst, '.')
	dst = strconv.AppendUint(dst, v.n, 10)
	dst = append(dst, '.')
	dst = strconv.AppendUint(dst, v.p, 10)
	if len(v.pre) > 0 {
		dst = append(dst, '-')
		dst = append(dst, v.pre...)
	}
	if len(v.meta) > 0 {
		dst = append(dst, '+')
		dst = append(dst, v.meta...)
	}
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
