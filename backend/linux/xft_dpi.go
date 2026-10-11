//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"strconv"
	"strings"
)

// xftDPI reads Xft.dpi from the resource database on the root window, which
// is where a desktop's display-scaling setting ends up.
func (x *Driver) xftDPI() float64 {
	value, ok := x.property(x.root, xResourceManager, atomString, 16384)
	if !ok {
		return 0
	}
	for line := range strings.Lines(string(value)) {
		name, setting, found := strings.Cut(line, ":")
		if !found || strings.TrimSpace(name) != "Xft.dpi" {
			continue
		}
		dpi, err := strconv.ParseFloat(strings.TrimSpace(setting), 64)
		if err != nil || dpi < 48 || dpi > 480 {
			return 0
		}
		return dpi
	}
	return 0
}
