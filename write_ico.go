package antui

import (
	"bytes"
	"encoding/binary"
	"github.com/gabrielluizsf/antui/canvas"
	"io"
)

// WriteICO writes the pictures as a Windows .ico file, largest first.
//
// It is worth having on its own — an .ico is what a file browser or an
// installer wants — and it is half of what the resource below is made of.
func WriteICO(w io.Writer, images ...*canvas.Canvas) error {
	pictures, err := iconPictures(images)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	binary.Write(&out, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(&out, binary.LittleEndian, uint16(1)) // an icon, not a cursor
	binary.Write(&out, binary.LittleEndian, uint16(len(pictures)))

	at := uint32(6 + 16*len(pictures))
	for _, picture := range pictures {
		writeIconEntry(&out, picture, func(size uint32) {
			binary.Write(&out, binary.LittleEndian, size)
			binary.Write(&out, binary.LittleEndian, at)
		})
		at += uint32(len(picture.data))
	}
	for _, picture := range pictures {
		out.Write(picture.data)
	}
	_, err = w.Write(out.Bytes())
	return err
}
