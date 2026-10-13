//go:build windows

package antui

import (
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
	procCreateWaitableTimerW   = kernel32.NewProc("CreateWaitableTimerW")
	procSetWaitableTimer       = kernel32.NewProc("SetWaitableTimer")
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
	if timer == 0 {
		timer, _, _ = procCreateWaitableTimerW.Call(0, 0)
	}
	if timer == 0 {
		time.Sleep(d)
		return
	}
	defer procCloseHandle.Call(timer)
	if !setWaitableTimer(timer, d) {
		time.Sleep(d)
		return
	}
	procWaitForSingleObject.Call(timer, waitInfinite)
}
