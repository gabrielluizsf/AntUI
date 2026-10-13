//go:build windows

package windows

import "syscall"

func utf16Ptr(s string) *uint16 {
	p, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		// A string with an interior NUL cannot be a title; an empty one can.
		p, _ = syscall.UTF16PtrFromString("")
	}
	return p
}
