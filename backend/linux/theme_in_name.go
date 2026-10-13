//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "strings"

// themeInName reads a theme's name, which is where a desktop spells its side
// out: a theme that is dark says so in its name or in the variant after a
// colon (Adwaita:dark), and turning dark on is exactly what that marker is
// for — so a name carrying neither is a light theme. No name at all is not
// an answer.
func themeInName(value string) (dark, ok bool) {
	if strings.TrimSpace(value) == "" {
		return false, false
	}
	return strings.Contains(strings.ToLower(value), "dark"), true
}
