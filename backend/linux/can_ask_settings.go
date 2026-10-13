//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "strings"

// canAskSettings reports whether there is a tool worth asking, and a desktop
// worth asking it about. The tool reads GNOME's settings, which is what a
// GNOME-family desktop answers through and what its desktops grew from; a
// desktop of another family keeps its answer in a file themeInConfig reads,
// and a session that names no desktop is a window manager with no theme of
// its own to report.
func canAskSettings(env func(string) string) bool {
	if settingsHelper() == "" {
		return false
	}
	names := strings.ToLower(env("XDG_CURRENT_DESKTOP") + ";" +
		env("XDG_SESSION_DESKTOP") + ";" + env("DESKTOP_SESSION"))
	for _, family := range []string{"gnome", "ubuntu", "cinnamon", "budgie", "unity", "pantheon"} {
		if strings.Contains(names, family) {
			return true
		}
	}
	return false
}
