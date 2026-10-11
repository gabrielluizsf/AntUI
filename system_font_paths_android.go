//go:build android

package antui

// systemFontPaths is the fallback, for a device whose configuration file
// cannot be read or does not parse.
func systemFontPaths() []string {
	return []string{
		"/system/fonts/Roboto-Regular.ttf",
		"/system/fonts/RobotoStatic-Regular.ttf",
		"/system/fonts/NotoSans-Regular.ttf",
		"/system/fonts/DroidSans.ttf",
	}
}
