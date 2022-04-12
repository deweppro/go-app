package error

import "github.com/deweppro/go-errors"

//nolint:golint
var (
	ErrDepRunning     = errors.New("dependencies is already running")
	ErrDepNotRunning  = errors.New("dependencies are not running yet")
	ErrDepEmpty       = errors.New("dependencies is empty")
	ErrServiceUnknown = errors.New("unknown service")
	ErrConfigUnknown  = errors.New("unknown config type")
	ErrBadFileFormat  = errors.New("is not a supported file format")
)
