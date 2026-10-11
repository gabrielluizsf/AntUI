//go:build windows

package windows

import (
	"syscall"
	"unsafe"
)

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
