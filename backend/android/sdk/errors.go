package sdk

import "errors"

// as is errors.As, kept behind a name so the doctor reads as prose. It is
// also safe on a nil error, which errors.As is not documented to be.
func as[T error](err error, target *T) bool {
	if err == nil {
		return false
	}
	return errors.As(err, target)
}
