//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os"
	"path/filepath"
)

// themeInConfig reads the files a desktop writes when it starts, in the
// order their own desktops are the likeliest ones to be running. A file that
// is not there is skipped rather than an error, and the first file with
// something to say wins.
//
// Only the two families that keep a theme in a file are read. A Qt outside
// of KDE has no such file — the desktops that would have one write these
// instead — and a desktop that keeps its answer somewhere else (an XML
// settings tree, a database) has said nothing to a program that reads files.
func themeInConfig(dir string) (dark, ok bool) {
	if dir == "" {
		return false, false
	}
	for _, probe := range []struct {
		path string
		read func(string) (bool, bool)
	}{
		{filepath.Join(dir, "kdeglobals"), kdeTheme},
		{filepath.Join(dir, "gtk-3.0", "settings.ini"), gtkTheme},
		{filepath.Join(dir, "gtk-4.0", "settings.ini"), gtkTheme},
	} {
		data, err := os.ReadFile(probe.path)
		if err != nil {
			continue
		}
		if dark, said := probe.read(string(data)); said {
			return dark, true
		}
	}
	return false, false
}
