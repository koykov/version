package version

import (
	"io"
	"strconv"

	"github.com/koykov/byteconv"
)

type Semver struct {
	m, n, p  uint32
	pr, meta []byte
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
	_ = ver
	// todo implement me
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
	v.pr = value
	return v
}

func (v *Semver) SetPreReleaseString(value string) *Semver {
	v.pr = byteconv.S2B(value)
	return v
}

func (v *Semver) SetMeta(value []byte) *Semver {
	v.meta = value
	return v
}

func (v *Semver) SetMetaString(value string) *Semver {
	v.meta = byteconv.S2B(value)
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
	return v.pr
}

func (v *Semver) PreReleaseString() string {
	return byteconv.B2S(v.pr)
}

func (v *Semver) Meta() []byte {
	return v.meta
}

func (v *Semver) MetaString() string {
	return byteconv.B2S(v.meta)
}

func (v *Semver) AppendBytes(dst []byte) []byte {
	dst = strconv.AppendUint(dst, uint64(v.m), 10)
	dst = append(dst, '.')
	dst = strconv.AppendUint(dst, uint64(v.n), 10)
	if v.p > 0 || len(v.pr) > 0 || len(v.meta) > 0 {
		dst = append(dst, '.')
		dst = strconv.AppendUint(dst, uint64(v.p), 10)
	}
	if len(v.pr) > 0 {
		dst = append(dst, '-')
		dst = append(dst, v.pr...)
	}
	if len(v.meta) > 0 {
		dst = append(dst, '+')
		dst = append(dst, v.meta...)
	}
	return dst
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
