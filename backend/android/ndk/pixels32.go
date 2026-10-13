//go:build android

package ndk

import "unsafe"

// Pixels32 is the buffer as 32-bit words, Stride*Height of them. It is the
// compositor's memory, not Go's, and stops being valid the moment the window
// is unlocked.
func (b Buffer) Pixels32() []uint32 {
	if b.Bits == nil || b.Format.Bytes() != 4 {
		return nil
	}
	return unsafe.Slice((*uint32)(b.Bits), b.Stride*b.Height)
}
