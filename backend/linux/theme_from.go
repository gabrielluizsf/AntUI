//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// themeFrom reads the color scheme from the sources that could know it,
// best first: the environment is set for this one program and is the most
// specific word about it, the desktop's settings are what the user chose for
// everything they run, the files a desktop writes at login are what is left
// when there is no tool to ask, and the colours a terminal exports describe
// one window rather than the desktop behind it. The first source that says
// something wins, and when none does the answer is that none did.
func themeFrom(env func(string) string, dir string, setting darkAnswer) (dark, ok bool) {
	if dark, said := themeInName(env("GTK_THEME")); said {
		return dark, true
	}
	if setting.ok {
		return setting.dark, true
	}
	if dark, said := themeInConfig(dir); said {
		return dark, true
	}
	return themeInColorFGBG(env("COLORFGBG"))
}
