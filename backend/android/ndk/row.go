//go:build android

package ndk

import (
	"unsafe"
)

// Row is one row of the buffer as bytes, whatever the format.
func (b Buffer) Row(y int) []byte {
	if b.Bits == nil || y < 0 || y >= b.Height {
		return nil
	}
	px := b.Format.Bytes()
	start := y * b.Stride * px
	return unsafe.Slice((*byte)(b.Bits), b.Stride*b.Height*px)[start : start+b.Width*px]
}
