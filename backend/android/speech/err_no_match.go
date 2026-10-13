//go:build android

package speech

import "errors"

// The reasons recognition stops, translated out of the engine's numbers.
var (
	ErrNoMatch      = errors.New("speech: nothing recognisable was said")
	ErrTimeout      = errors.New("speech: nothing was said")
	ErrNoNetwork    = errors.New("speech: recognition needs a network and there is none")
	ErrBusy         = errors.New("speech: the recogniser is busy")
	ErrNoPermission = errors.New("speech: this app has not been granted RECORD_AUDIO")
)
