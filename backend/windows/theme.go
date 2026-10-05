//go:build windows

package windows

import (
	"syscall"
	"unsafe"
)

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

// hkeyCurrentUser is HKEY_CURRENT_USER: 0x80000001 with the rest of the
// machine word filled in, which is how the windows.h header gets there — it
// casts a LONG — and what the registry expects to be handed. keyRead is the
// access the theme value is read with, and regDword the type it is kept in.
const (
	hkeyCurrentUser = ^uintptr(0x7FFFFFFE) // 0x80000001 sign-extended
	keyRead         = 0x20019
	regDword        = 4
)

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

// touchTheme marks what the registry last said stale, so the next frame asks
// it again. The window procedure calls it on the setting change Windows
// posts when the theme is flipped.
func (d *Driver) touchTheme() { d.themeStale = true }

// readAppsTheme reads the one value out of the user's own hive. A key the
// account does not have, a value that is not there, or a value that is not a
// DWORD all come back as no answer rather than as a theme.
func readAppsTheme() (dark, known bool) {
	subKey, err := syscall.UTF16PtrFromString(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`)
	if err != nil {
		return false, false
	}
	var hkey uintptr
	opened, _, _ := procRegOpenKeyExW.Call(
		hkeyCurrentUser,
		uintptr(unsafe.Pointer(subKey)),
		0,
		keyRead,
		uintptr(unsafe.Pointer(&hkey)),
	)
	if opened != 0 {
		return false, false
	}
	defer procRegCloseKey.Call(hkey)

	name, err := syscall.UTF16PtrFromString("AppsUseLightTheme")
	if err != nil {
		return false, false
	}
	var value, typ, size uint32 = 0, 0, 4
	read, _, _ := procRegQueryValueExW.Call(
		hkey,
		uintptr(unsafe.Pointer(name)),
		0,
		uintptr(unsafe.Pointer(&typ)),
		uintptr(unsafe.Pointer(&value)),
		uintptr(unsafe.Pointer(&size)),
	)
	if read != 0 || typ != regDword || size < 4 {
		return false, false
	}
	// AppsUseLightTheme is 1 while the light theme is on, so dark is its
	// zero.
	return value == 0, true
}
