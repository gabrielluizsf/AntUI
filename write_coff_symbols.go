package antui

import (
	"bytes"
	"encoding/binary"
)

// writeCOFFSymbols writes the symbol table: the section, its auxiliary
// record, and the empty string table.
func writeCOFFSymbols(out *bytes.Buffer, section resourceSection) {
	// Eight characters is exactly what ".rsrc" fits in, so there is no
	// string table to write.
	out.Write([]byte{'.', 'r', 's', 'r', 'c', 0, 0, 0})
	binary.Write(out, binary.LittleEndian, uint32(0)) // value: the start of the section
	binary.Write(out, binary.LittleEndian, uint16(1)) // section number, from one
	binary.Write(out, binary.LittleEndian, uint16(0)) // type: none
	out.WriteByte(3)                                  // storage class: static
	out.WriteByte(1)                                  // one auxiliary record

	binary.Write(out, binary.LittleEndian, uint32(len(section.data)))
	binary.Write(out, binary.LittleEndian, uint16(len(section.fixups)))
	binary.Write(out, binary.LittleEndian, uint16(0)) // line numbers
	binary.Write(out, binary.LittleEndian, uint32(0)) // checksum
	binary.Write(out, binary.LittleEndian, uint16(0)) // number of the section this is a copy of
	// The selection byte and three unused ones. An auxiliary record is
	// eighteen bytes exactly, the same as a symbol, and one byte short of
	// that makes every reader lose the string table.
	out.Write([]byte{0, 0, 0, 0})

	// An empty string table, which is four bytes saying its own size. A COFF
	// file without one is a COFF file some readers refuse.
	binary.Write(out, binary.LittleEndian, uint32(4))
}
