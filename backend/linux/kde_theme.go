//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"strings"
)

// kdeTheme reads the scheme KDE names in kdeglobals: ColorScheme is the
// theme's name and LookAndFeelPackage is the one a global theme writes
// instead, and both are names, which is how a desktop says its side.
func kdeTheme(text string) (dark, ok bool) {
	for line := range strings.Lines(text) {
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		switch strings.TrimSpace(key) {
		case "ColorScheme", "LookAndFeelPackage":
			return themeInName(value)
		}
	}
	return false, false
}
