package version

import (
	"strings"
	"testing"
)

type tcSV struct {
	src, pre, meta string
	m, n, p        uint64
	err            error
}

var tcSMs = []tcSV{
	{src: "     v1.0.0", m: 1},
	{src: "     4.3.45    ", m: 4, n: 3, p: 45},
	{src: "1.2.3", m: 1, n: 2, p: 3},
	{src: "1.2.3-alpha.01", m: 1, n: 2, p: 3, pre: "alpha.01"},
	{src: "1.2.3+test.01", m: 1, n: 2, p: 3, meta: "test.01"},
	{src: "1.2.3-alpha.-1", m: 1, n: 2, p: 3, pre: "alpha.-1"},
	{src: "v1.2.3", m: 1, n: 2, p: 3},
	{src: "1.0", err: ErrBadSemver},
	{src: "v1.0", err: ErrBadSemver},
	{src: "1", err: ErrBadSemver},
	{src: "v1", err: ErrBadSemver},
	{src: "1.2.beta", err: ErrBadSemver},
	{src: "v1.2.beta", err: ErrBadSemver},
	{src: "foo", err: ErrBadSemver},
	{src: "1.2-5", err: ErrBadSemver},
	{src: "v1.2-5", err: ErrBadSemver},
	{src: "1.2-beta.5", err: ErrBadSemver},
	{src: "v1.2-beta.5", err: ErrBadSemver},
	{src: "\n1.2", err: ErrBadSemver},
	{src: "\nv1.2", err: ErrBadSemver},
	{src: "1.2.0-x.Y.0+metadata", m: 1, n: 2, pre: "x.Y.0", meta: "metadata"},
	{src: "v1.2.0-x.Y.0+metadata", m: 1, n: 2, pre: "x.Y.0", meta: "metadata"},
	{src: "1.2.0-x.Y.0+metadata-width-hypen", m: 1, n: 2, pre: "x.Y.0", meta: "metadata-width-hypen"},
	{src: "v1.2.0-x.Y.0+metadata-width-hypen", m: 1, n: 2, pre: "x.Y.0", meta: "metadata-width-hypen"},
	{src: "1.2.3-rc1-with-hypen", m: 1, n: 2, p: 3, pre: "rc1-with-hypen"},
	{src: "v1.2.3-rc1-with-hypen", m: 1, n: 2, p: 3, pre: "rc1-with-hypen"},
	{src: "1.2.3.4", err: ErrBadSemver},
	{src: "v1.2.3.4", err: ErrBadSemver},
	{src: "1.2.2147483648", m: 1, n: 2, p: 2147483648},
	{src: "1.2147483648.3", m: 1, n: 2147483648, p: 3},
	{src: "2147483648.3.0", m: 2147483648, n: 3},
}

func TestSemverParse(t *testing.T) {
	for _, stg := range tcSMs {
		t.Run(stg.src, func(t *testing.T) {
			var ver Semver
			err := ver.ParseString(stg.src)
			if err != nil && stg.err == nil {
				t.Error(err)
				return
			}
		})
	}
}

func TestSemverMarshal(t *testing.T) {
	for _, stg := range tcSMs {
		t.Run(stg.src, func(t *testing.T) {
			if stg.err != nil {
				return
			}
			ver := NewSemver(stg.m, stg.n, stg.p, stg.pre, stg.meta)
			s := ver.String()
			if i := strings.Index(stg.src, s); i == -1 {
				t.FailNow()
			}
		})
	}
}

func BenchmarkSemverParse(b *testing.B) {
	for _, stg := range tcSMs {
		b.Run(stg.src, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = ParseSemverString(stg.src)
			}
		})
	}
}
