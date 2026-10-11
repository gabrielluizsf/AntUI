//go:build windows

package windows

type trackMouseEvent struct {
	Size      uint32
	Flags     uint32
	HWndTrack uintptr
	HoverTime uint32
}
