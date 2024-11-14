package version

import "errors"

var (
	ErrEmpty   = errors.New("empty version")
	ErrShort   = errors.New("data is too short")
	ErrBinLen8 = errors.New("binary data length less than 8 bytes")
)
