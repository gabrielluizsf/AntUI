//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// themeInSettings asks the desktop's settings tool which scheme it is
// drawing in.
//
// The two keys are the ones GNOME grew and its desktops have kept writing
// into ever since — Cinnamon and MATE both write the theme back through
// them, which is why GTK applications pick a desktop's choice up at all. The
// scheme is asked first, because it is the answer for the newer desktops
// that have one; "default" there means the theme already is what the user
// asked for, and the theme's name is then where that desktop says which one
// that is.
func themeInSettings(ask settingsAsk) (dark, ok bool) {
	switch scheme := unquote(ask("get", "org.gnome.desktop.interface", "color-scheme")); scheme {
	case "prefer-dark":
		return true, true
	case "prefer-light":
		return false, true
	}
	return themeInName(unquote(ask("get", "org.gnome.desktop.interface", "gtk-theme")))
}
