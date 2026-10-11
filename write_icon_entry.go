package antui

import (
	"bytes"
	"encoding/binary"
)

// writeIconEntry writes the sixteen-byte description of one picture that both
// an .ico and a group resource begin with, and calls `tail` for the four
// bytes at the end, which is a file offset in one and a resource number in
// the other.
func writeIconEntry(out *bytes.Buffer, picture iconPicture, tail func(size uint32)) {
	// 256 is written as nought: the field is one byte.
	out.WriteByte(byte(picture.width % 256))
	out.WriteByte(byte(picture.height % 256))
	out.WriteByte(0)                                   // colours in the palette: none, it is a true-colour icon
	out.WriteByte(0)                                   // reserved
	binary.Write(out, binary.LittleEndian, uint16(1))  // planes
	binary.Write(out, binary.LittleEndian, uint16(32)) // bits a pixel
	tail(uint32(len(picture.data)))
}
