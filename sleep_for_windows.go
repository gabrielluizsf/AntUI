//go:build windows

package antui

import (
	"syscall"
	"time"
)

const (
	createWaitableTimerHighResolution = 0x00000002
	timerModifyState                  = 0x0002
	synchronize                       = 0x00100000
	waitInfinite                      = 0xffffffff
)

var (
	procCreateWaitableTimerExW = kernel32.NewProc("CreateWaitableTimerExW")
	procSetWaitableTimer       = kernel32.NewProc("SetWaitableTimer")

	// Windows before 10 1803 has no high-resolution waitable-timer flag. Its
	// regular timer can still be made precise enough for a 60 Hz frame by
	// temporarily requesting a 1 ms timer resolution.
	winmm               = syscall.NewLazyDLL("winmm.dll")
	procTimeBeginPeriod = winmm.NewProc("timeBeginPeriod")
	procTimeEndPeriod   = winmm.NewProc("timeEndPeriod")
)

// sleepFor uses a relative high-resolution waitable timer. Unlike Sleep, this
// does not depend on the system timer tick, which is commonly too coarse for
// a 16.666 ms frame. The regular timer remains a fallback for older Windows
// versions that do not support the high-resolution flag.
func sleepFor(d time.Duration) {
	if d <= 0 {
		return
	}

	timer, _, _ := procCreateWaitableTimerExW.Call(
		0,
		0,
		createWaitableTimerHighResolution,
		timerModifyState|synchronize,
	)
	if timer != 0 {
		defer procCloseHandle.Call(timer)
		if setWaitableTimer(timer, d) {
			procWaitForSingleObject.Call(timer, waitInfinite)
			return
		}
	}

	// The high-resolution flag is unavailable on older Windows versions and
	// may also be unavailable under compatibility layers. Do not fall back to
	// a bare time.Sleep: that is the path that turns 16.6 ms into ~20 ms.
	procTimeBeginPeriod.Call(1)
	time.Sleep(d)
	procTimeEndPeriod.Call(1)
}
