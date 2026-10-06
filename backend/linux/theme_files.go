//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os"
	"path/filepath"
	"strings"
)

// configDir is where a desktop keeps the settings it writes at login:
// XDG_CONFIG_HOME when the session named one, the home directory's .config
// when it did not, and nothing when there is neither — a session with no
// home directory has no settings file to read, and says so by having none.
func configDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return dir
}

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
