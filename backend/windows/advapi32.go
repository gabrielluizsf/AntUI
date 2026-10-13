//go:build windows

package windows

import "syscall"

// The color scheme Windows paints its own apps in is a registry value rather
// than a call with a window in it: AppsUseLightTheme, under the Personalize
// key, which the user flips in Settings. Windows posts a setting change when
// that happens, so the value is read once, kept, and read again only when the
// message says it moved — never once per frame, which would be a registry
// lookup on every frame the window draws.
var (
	advapi32             = syscall.NewLazyDLL("advapi32.dll")
	procRegOpenKeyExW    = advapi32.NewProc("RegOpenKeyExW")
	procRegQueryValueExW = advapi32.NewProc("RegQueryValueExW")
	procRegCloseKey      = advapi32.NewProc("RegCloseKey")
)
