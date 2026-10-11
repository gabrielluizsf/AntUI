//go:build windows

package windows

import (
	"syscall"
	"unsafe"
)

// utf16PtrToString reads a zero-terminated UTF-16 string out of memory the
// clipboard owns. It stops at a length no clipboard text reaches, so a
// handle that is not what it claims cannot walk off the end of the world.
//
// The pointer comes from GlobalLock rather than from Go's heap, so it is
// walked as one — `unsafe.Slice` over the locked block, which is what the
// checker wants to see rather than arithmetic on a uintptr.
func utf16PtrToString(at uintptr) string {
	const most = 1 << 20
	block := unsafe.Slice((*uint16)(unsafe.Pointer(at)), most)
	for i, c := range block {
		if c == 0 {
			return syscall.UTF16ToString(block[:i])
		}
	}
	return syscall.UTF16ToString(block)
}
