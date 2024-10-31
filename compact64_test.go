package version

import "testing"

type tc64 struct {
	raw        string
	m, n, p, r uint16
}

var tcs64 = []tc64{
	{"", 0, 0, 0, 0},
	{"0", 0, 0, 0, 0},
	{"1", 1, 0, 0, 0},
	{"1.0", 1, 0, 0, 0},
	{"1.0.1", 1, 0, 1, 0},
	{"1.0.1.7", 1, 0, 1, 7},
	{"5.12.134", 5, 12, 134, 0},
}

func TestCompact64Parse(t *testing.T) {
	for _, c := range tcs64 {
		t.Run(c.raw, func(t *testing.T) {
			ver, err := ParseCompact64String(c.raw)
			if err != nil && err != ErrEmpty {
				t.Error(err)
			}
			if ver.Major() != c.m || ver.Minor() != c.n || ver.Patch() != c.p || ver.Revision() != c.r {
				t.FailNow()
			}
		})
	}
}

func BenchmarkCompact64Parse(b *testing.B) {
	for _, c := range tcs64 {
		b.Run(c.raw, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = ParseCompact64String(c.raw)
			}
		})
	}
}
