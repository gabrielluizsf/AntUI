//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"strings"
)

// gtkTheme reads a GTK settings.ini: the theme's name, and the flag that
// asks for dark whatever that name is. The flag on wins over the name,
// since it is the one that turns a light theme dark; the name wins over the
// flag off, which nearly every settings.ini carries and which only says the
// theme was not overridden.
func gtkTheme(text string) (dark, ok bool) {
	var name string
	forced, released := false, false
	for line := range strings.Lines(text) {
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		switch key {
		case "gtk-application-prefer-dark-theme":
			switch strings.ToLower(value) {
			case "1", "true", "yes":
				forced = true
			case "0", "false", "no":
				released = true
			}
		case "gtk-theme-name":
			name = value
		}
	}
	if forced {
		return true, true
	}
	if dark, said := themeInName(name); said {
		return dark, true
	}
	if released {
		return false, true
	}
	return false, false
}
