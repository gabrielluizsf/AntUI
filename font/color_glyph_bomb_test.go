package font

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/png"
	"testing"
)

// pngClaiming is a PNG of a few dozen bytes whose header says it is w by h,
// with one empty block of pixel data after it.
func pngClaiming(w, h uint32) []byte {
	var buf bytes.Buffer
	buf.WriteString("\x89PNG\r\n\x1a\n")
	chunk := func(typ string, data []byte) {
		binary.Write(&buf, binary.BigEndian, uint32(len(data)))
		c := append([]byte(typ), data...)
		buf.Write(c)
		binary.Write(&buf, binary.BigEndian, crc32.ChecksumIEEE(c))
	}
	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr, w)
	binary.BigEndian.PutUint32(ihdr[4:], h)
	ihdr[8], ihdr[9] = 8, 6
	chunk("IHDR", ihdr)
	chunk("IDAT", []byte{0x78, 0x9c, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01})
	chunk("IEND", nil)
	return buf.Bytes()
}

// onePictureFont is a font whose one colour glyph is the PNG given, written as
// image format 19: its length, then the picture.
func onePictureFont(pic []byte) *TTF {
	cbdt := binary.BigEndian.AppendUint32(nil, uint32(len(pic)))
	cbdt = append(cbdt, pic...)
	return &TTF{cbdt: cbdt, colorStrikes: []colorStrike{{ppem: 16, glyphs: map[int]colorBitmap{
		0: {off: 0, length: len(cbdt), imageFormat: 19, inIndex: true},
	}}}}
}

// TestColorGlyphRefusesADecompressionBomb is an emoji font whose picture for
// a glyph says it is sixty thousand pixels on a side: the PNG decoder would
// allocate fourteen gigabytes for it before reading a pixel.
func TestColorGlyphRefusesADecompressionBomb(t *testing.T) {
	if cg := onePictureFont(pngClaiming(60000, 60000)).ColorGlyph(0, 16); cg != nil {
		t.Fatal("the bomb was decoded, want the glyph left without its picture")
	}
}

// TestColorGlyphStillReadsAnEmoji keeps the limit above what an emoji font
// cuts: a picture of a hundred and thirty-six pixels is read.
func TestColorGlyphStillReadsAnEmoji(t *testing.T) {
	var pic bytes.Buffer
	if err := png.Encode(&pic, image.NewNRGBA(image.Rect(0, 0, 136, 128))); err != nil {
		t.Fatal(err)
	}
	if cg := onePictureFont(pic.Bytes()).ColorGlyph(0, 16); cg == nil {
		t.Fatal("an emoji-sized picture was refused")
	}
}
