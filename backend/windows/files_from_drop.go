//go:build windows

package windows

import (
	"syscall"
	"unsafe"
)

// filesFromDrop is the paths an HDROP holds, whether it came from a drag or
// off the clipboard.
func filesFromDrop(hdrop uintptr) []string {
	count, _, _ := procDragQueryFileW.Call(hdrop, dragQueryCount, 0, 0)
	if count == 0 || count > 4096 {
		return nil
	}
	out := make([]string, 0, count)
	for i := uintptr(0); i < count; i++ {
		// Asked for the length first, because a path may be long and the old
		// 260-character limit is not one any more.
		length, _, _ := procDragQueryFileW.Call(hdrop, i, 0, 0)
		if length == 0 {
			continue
		}
		buf := make([]uint16, length+1)
		got, _, _ := procDragQueryFileW.Call(hdrop, i,
			uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if got == 0 {
			continue
		}
		out = append(out, syscall.UTF16ToString(buf[:got]))
	}
	return out
}
