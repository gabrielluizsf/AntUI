//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"strconv"
	"strings"
)

// themeInColorFGBG reads the colour pair a terminal exports in COLORFGBG,
// "foreground;background" as palette indices. Black through cyan is a dark
// background and white or bright white is a light one; the indices in
// between are a pair someone found readable, not a statement about the
// desktop, and say nothing.
func themeInColorFGBG(value string) (dark, ok bool) {
	_, background, found := strings.Cut(value, ";")
	if !found {
		return false, false
	}
	index, err := strconv.Atoi(strings.TrimSpace(background))
	if err != nil {
		return false, false
	}
	switch {
	case index >= 0 && index <= 6:
		return true, true
	case index == 7 || index == 15:
		return false, true
	}
	return false, false
}
