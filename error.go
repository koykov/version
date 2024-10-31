package version

import "errors"

var (
	ErrEmpty   = errors.New("empty version")
	ErrBinLen8 = errors.New("binary data length less than 8 bytes")
)
