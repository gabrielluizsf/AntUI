//go:build android

package speech

import (
	"fmt"
)

// recogniserError turns the engine's number into something a person could
// act on.
func recogniserError(code int32) error {
	switch code {
	case 1, 2:
		return ErrNoNetwork
	case 6:
		return ErrTimeout
	case 7:
		return ErrNoMatch
	case 8:
		return ErrBusy
	case 9:
		return ErrNoPermission
	}
	return fmt.Errorf("speech: the recogniser stopped (%d)", code)
}
