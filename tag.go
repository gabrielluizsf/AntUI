//go:build linux && !android

package antui

import (
	"fmt"
)

func (r *read) tag(want byte) error {
	if r.at >= len(r.b) {
		return errShort
	}
	if got := r.b[r.at]; got != want {
		return fmt.Errorf("antui: the server sent a %q where a %q belongs", got, want)
	}
	r.at++
	return nil
}
