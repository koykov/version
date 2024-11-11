package version

import "testing"

type tcSV struct {
	src, pr, meta string
	m, n, p       uint32
	err           error
}

var tcSMs = []tcSV{
	{
		src: "     v1.0",
		m:   1,
	},
	{
		src: "     4.3.45    ",
		m:   4,
		n:   3,
		p:   45,
	},
	{src: "1.2.3"},
	{src: "1.2.3-alpha.01", err: ErrBadSemver},
	{src: "1.2.3+test.01"},
	{src: "1.2.3-alpha.-1"},
	{src: "v1.2.3", err: ErrBadSemver},
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
	{src: "1.2.0-x.Y.0+metadata"},
	{src: "v1.2.0-x.Y.0+metadata", err: ErrBadSemver},
	{src: "1.2.0-x.Y.0+metadata-width-hypen"},
	{src: "v1.2.0-x.Y.0+metadata-width-hypen", err: ErrBadSemver},
	{src: "1.2.3-rc1-with-hypen"},
	{src: "v1.2.3-rc1-with-hypen", err: ErrBadSemver},
	{src: "1.2.3.4", err: ErrBadSemver},
	{src: "v1.2.3.4", err: ErrBadSemver},
	{src: "1.2.2147483648"},
	{src: "1.2147483648.3"},
	{src: "2147483648.3.0"},
}

func TestSemverParse(t *testing.T) {
	for _, stg := range tcSMs {
		t.Run(stg.src, func(t *testing.T) {
			var ver Semver
			_ = ver.ParseString(stg.src)
		})
	}
}
