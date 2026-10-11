//go:build windows

package antui

import (
	"fmt"
)

// failed turns an HRESULT into an error. Anything with the top bit set is a
// failure; everything else, S_FALSE included, is not.
func failed(hr uintptr, what string) error {
	if int32(hr) >= 0 {
		return nil
	}
	if uint32(hr) == errNotFound {
		return fmt.Errorf("%w: %s found no device", ErrNoAudio, what)
	}
	return fmt.Errorf("antui: %s failed (0x%08X)", what, uint32(hr))
}
