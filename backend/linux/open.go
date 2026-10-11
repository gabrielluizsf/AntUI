//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend"
)

// Open connects to the display and creates the window that Options describes.
func Open(win backend.Face, opts backend.Options) (*Driver, error) {
	if win == nil {
		return nil, errors.New("antui: x11 backend with no window to drive")
	}
	d := &Driver{alive: true}
	if err := d.open(win, opts.Title, opts.Width, opts.Height); err != nil {
		return nil, err
	}
	return d, nil
}
