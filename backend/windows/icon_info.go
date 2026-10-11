//go:build windows

package windows

type iconInfo struct {
	icon     int32 // 1 for an icon, 0 for a cursor
	xHotspot uint32
	yHotspot uint32
	mask     uintptr
	colour   uintptr
}
