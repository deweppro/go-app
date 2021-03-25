package app

import (
	"github.com/pkg/errors"
)

//nolint:golint
var (
	ErrDepRunning    = errors.New("dependencies is already running")
	ErrDepNotRunning = errors.New("dependencies are not running yet")
	ErrDepEmpty      = errors.New("dependencies is empty")
	ErrDepUnknown    = errors.New("unknown dependency")
	ErrBadAction     = errors.New("is not a supported action")
	ErrBadFileFormat = errors.New("is not a supported file format")
)

//WrapErrors combining multiple errors
func WrapErrors(err1, err2 error, message string) error {
	if err2 == nil {
		return err1
	}
	if err1 == nil {
		return errors.Wrap(err2, message)
	}
	return errors.Wrap(err1, errors.Wrap(err2, message).Error())
}
