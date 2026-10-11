//go:build windows

package windows

import (
	"unsafe"
)

// DisplayRefresh is how often the display repaints, in hertz.
func (d *Driver) DisplayRefresh() int {
	var mode [devModeSize]byte
	*(*uint16)(unsafe.Pointer(&mode[devModeSizeOffset])) = devModeSize

	ok, _, _ := procEnumDisplaySettingsW.Call(0, enumCurrentSettings,
		uintptr(unsafe.Pointer(&mode[0])))
	if ok == 0 {
		return 0
	}
	hz := *(*uint32)(unsafe.Pointer(&mode[devModeFrequencyOff]))
	// 0 and 1 both mean "the hardware default", which is not a number.
	if hz <= 1 {
		return 0
	}
	return int(hz)
}
