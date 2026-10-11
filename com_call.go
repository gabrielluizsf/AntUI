//go:build windows

package antui

import (
	"syscall"
	"unsafe"
)

// call invokes the method at an index of an interface's vtable. `this` goes
// in front of the arguments, which is what makes a COM call a COM call.
func comCall(this uintptr, index int, args ...uintptr) uintptr {
	vtable := *(**[64]uintptr)(unsafe.Pointer(this))
	all := make([]uintptr, 0, len(args)+1)
	all = append(all, this)
	all = append(all, args...)
	ret, _, _ := syscall.SyscallN(vtable[index], all...)
	return ret
}
