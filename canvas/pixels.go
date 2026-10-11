package canvas

import "unsafe"

// pixels is a run of colours as the plain numbers the assembly takes.
func pixels(run []Color) *uint32 {
	return (*uint32)(unsafe.Pointer(&run[0]))
}
