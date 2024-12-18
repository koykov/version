package version

import "testing"

type tc32 struct {
	raw        string
	m, n, p, r uint8
}

var tcs32 = []tc32{
	{"", 0, 0, 0, 0},
	{"0", 0, 0, 0, 0},
	{"1", 1, 0, 0, 0},
	{"1.0", 1, 0, 0, 0},
	{"1.0.1", 1, 0, 1, 0},
	{"1.0.1.7", 1, 0, 1, 7},
	{"5.12.134", 5, 12, 134, 0},
}

func TestCompact32Parse(t *testing.T) {
	for _, c := range tcs32 {
		t.Run(c.raw, func(t *testing.T) {
			ver, err := ParseVersion32String(c.raw)
			if err != nil && err != ErrEmpty {
				t.Error(err)
			}
			if ver.Major() != c.m || ver.Minor() != c.n || ver.Patch() != c.p || ver.Revision() != c.r {
				t.FailNow()
			}
		})
	}
}

func TestCompact32Marshal(t *testing.T) {
	for _, c := range tcs32 {
		if len(c.raw) < 3 {
			continue
		}
		t.Run(c.raw, func(t *testing.T) {
			var ver Version32
			ver.SetMajor(c.m).
				SetMinor(c.n).
				SetPatch(c.p).
				SetRevision(c.r)
			s := ver.String()
			if s != c.raw {
				t.Log(s)
				t.FailNow()
			}
		})
	}
}

func BenchmarkCompact32Parse(b *testing.B) {
	for _, c := range tcs32 {
		b.Run(c.raw, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = ParseVersion32String(c.raw)
			}
		})
	}
}

func BenchmarkCompact32Marshal(b *testing.B) {
	for _, c := range tcs32 {
		if len(c.raw) < 3 {
			continue
		}
		b.Run(c.raw, func(b *testing.B) {
			b.ReportAllocs()
			var buf []byte
			for i := 0; i < b.N; i++ {
				var ver Version32
				ver.SetMajor(c.m).
					SetMinor(c.n).
					SetPatch(c.p).
					SetRevision(c.r)
				buf = ver.AppendBytes(buf[:0])
			}
		})
	}
}
