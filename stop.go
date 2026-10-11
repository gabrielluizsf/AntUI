//go:build linux && !android

package antui

import (
	"fmt"
)

// stop records why the sound stopped, unless it stopped because it was told
// to.
func (p *pulse) stop(err error) {
	if p.audio.stopped() {
		return
	}
	p.audio.fail(fmt.Errorf("antui: the sound server went away: %w", err))
}
