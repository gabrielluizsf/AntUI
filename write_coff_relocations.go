package antui

import (
	"bytes"
	"encoding/binary"
)

// writeCOFFRelocations writes the fixups, all against the section's own
// symbol: the value already in the field is the offset from the start of the
// section, and the linker adds where the section ended up.
func writeCOFFRelocations(out *bytes.Buffer, section resourceSection, relocation uint16) {
	for _, at := range section.fixups {
		binary.Write(out, binary.LittleEndian, at)
		binary.Write(out, binary.LittleEndian, uint32(0)) // the first symbol
		binary.Write(out, binary.LittleEndian, relocation)
	}
}
