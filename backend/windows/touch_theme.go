//go:build windows

package windows

// touchTheme marks what the registry last said stale, so the next frame asks
// it again. The window procedure calls it on the setting change Windows
// posts when the theme is flipped.
func (d *Driver) touchTheme() { d.themeStale = true }
