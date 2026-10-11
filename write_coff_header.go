package antui

import (
	"bytes"
	"encoding/binary"
)

// writeCOFFHeader writes the COFF file header and the one section header.
func writeCOFFHeader(out *bytes.Buffer, section resourceSection, machine uint16) {
	const headerSize = 20
	const sectionHeaderSize = 40
	rawAt := uint32(headerSize + sectionHeaderSize)
	relocationsAt := rawAt + uint32(len(section.data))
	symbolsAt := relocationsAt + uint32(10*len(section.fixups))

	binary.Write(out, binary.LittleEndian, machine)
	binary.Write(out, binary.LittleEndian, uint16(1)) // one section
	binary.Write(out, binary.LittleEndian, uint32(0)) // no timestamp: same bytes every time
	binary.Write(out, binary.LittleEndian, symbolsAt)
	binary.Write(out, binary.LittleEndian, uint32(2)) // the section symbol and its aux
	binary.Write(out, binary.LittleEndian, uint16(0)) // no optional header
	binary.Write(out, binary.LittleEndian, uint16(0))

	out.Write([]byte{'.', 'r', 's', 'r', 'c', 0, 0, 0})
	binary.Write(out, binary.LittleEndian, uint32(0)) // virtual size: an object has none
	binary.Write(out, binary.LittleEndian, uint32(0)) // and no address yet
	binary.Write(out, binary.LittleEndian, uint32(len(section.data)))
	binary.Write(out, binary.LittleEndian, rawAt)
	binary.Write(out, binary.LittleEndian, relocationsAt)
	binary.Write(out, binary.LittleEndian, uint32(0)) // no line numbers
	binary.Write(out, binary.LittleEndian, uint16(len(section.fixups)))
	binary.Write(out, binary.LittleEndian, uint16(0))
	binary.Write(out, binary.LittleEndian, uint32(0x40000040)) // initialised, readable
}
