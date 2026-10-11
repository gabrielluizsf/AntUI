package antui

import (
	"bytes"
	"encoding/binary"
)

// groupIcon is the little directory that names the pictures. It is what
// Windows looks up first, and it is why the icons need numbers.
func groupIcon(pictures []iconPicture) []byte {
	var out bytes.Buffer
	binary.Write(&out, binary.LittleEndian, uint16(0))
	binary.Write(&out, binary.LittleEndian, uint16(1))
	binary.Write(&out, binary.LittleEndian, uint16(len(pictures)))
	for i, picture := range pictures {
		writeIconEntry(&out, picture, func(size uint32) {
			binary.Write(&out, binary.LittleEndian, size)
			// The number of the RT_ICON that holds it. They are numbered from
			// one, because nought is not a resource name.
			binary.Write(&out, binary.LittleEndian, uint16(i+1))
		})
	}
	return out.Bytes()
}
