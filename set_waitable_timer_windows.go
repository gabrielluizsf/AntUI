//go:build windows

package antui

import (
	"time"
	"unsafe"
)

func setWaitableTimer(timer uintptr, d time.Duration) bool {
	// Windows represents a relative due time as negative 100 ns units. Round up
	// so the frame is never released before its deadline.
	units := (d + 100*time.Nanosecond - 1) / (100 * time.Nanosecond)
	due := -int64(units)
	ok, _, _ := procSetWaitableTimer.Call(
		timer,
		uintptr(unsafe.Pointer(&due)),
		0,
		0,
		0,
		0,
	)
	return ok != 0
}
