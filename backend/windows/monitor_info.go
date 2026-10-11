//go:build windows

package windows

type monitorInfo struct {
	Size      uint32
	RcMonitor rect
	RcWork    rect
	Flags     uint32
}
