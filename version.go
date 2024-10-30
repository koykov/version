package version

import (
	"encoding"
	"fmt"
	"io"
)

type Version interface {
	fmt.Stringer
	io.WriterTo
	WriteBinaryTo(w io.Writer) (int64, error)
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}
