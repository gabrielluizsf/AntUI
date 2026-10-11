//go:build windows

package windows

// SystemDark is the scheme Windows paints its own interface in: dark when
// AppsUseLightTheme is off. known is false when the value is not there at
// all — a Windows that has never been asked has no answer, and half an
// answer would be a guess at a theme.
func (d *Driver) SystemDark() (dark, known bool) {
	if d == nil {
		return false, false
	}
	if d.themeRead && !d.themeStale {
		return d.themeDark, d.themeKnown
	}
	d.themeDark, d.themeKnown = readAppsTheme()
	d.themeRead, d.themeStale = true, false
	return d.themeDark, d.themeKnown
}
