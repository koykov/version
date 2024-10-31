package version

import (
	"encoding"
	"fmt"
	"io"
)

type Version interface {
	fmt.Stringer
	io.WriterTo
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	Parse(ver []byte) error
	ParseString(ver string) error
	WriteBinaryTo(w io.Writer) (int64, error)
	Bytes() []byte
}
