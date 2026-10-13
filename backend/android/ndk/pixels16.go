//go:build android

package ndk

import "unsafe"

// Pixels16 is the buffer as 16-bit words, for RGB565.
func (b Buffer) Pixels16() []uint16 {
	if b.Bits == nil || b.Format != RGB565 {
		return nil
	}
	return unsafe.Slice((*uint16)(b.Bits), b.Stride*b.Height)
}
