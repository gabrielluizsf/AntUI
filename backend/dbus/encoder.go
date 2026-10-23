package dbus

import "encoding/binary"

// encoder writes D-Bus's wire format, little-endian. Every value is aligned
// to its own size counted from the start of buf, which is why a message's
// header and its body each get an encoder of their own: both start on an
// 8-byte boundary of the message, so alignment from either start is the same
// as alignment from the message's.
type encoder struct{ buf []byte }

// align pads with zeros up to the next multiple of n.
func (e *encoder) align(n int) {
	for len(e.buf)%n != 0 {
		e.buf = append(e.buf, 0)
	}
}

func (e *encoder) byte(b byte) { e.buf = append(e.buf, b) }

func (e *encoder) uint32(v uint32) {
	e.align(4)
	e.buf = binary.LittleEndian.AppendUint32(e.buf, v)
}

func (e *encoder) int32(v int32) { e.uint32(uint32(v)) }

// string is a STRING or an OBJECT_PATH: a 32-bit length, the bytes and a
// NUL the length does not count.
func (e *encoder) string(s string) {
	e.uint32(uint32(len(s)))
	e.buf = append(e.buf, s...)
	e.buf = append(e.buf, 0)
}

// signature is a SIGNATURE: the same as a string with an 8-bit length.
func (e *encoder) signature(s string) {
	e.buf = append(e.buf, byte(len(s)))
	e.buf = append(e.buf, s...)
	e.buf = append(e.buf, 0)
}

// emptyArray is an array with no elements. The padding to the element's
// alignment is written even then, and not counted in the length.
func (e *encoder) emptyArray(elemAlign int) {
	e.uint32(0)
	e.align(elemAlign)
}
