package svg

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"image/gif"
	"strings"
	"testing"
)

// pngClaiming is a PNG of a few dozen bytes whose header says it is w by h,
// with one empty block of pixel data after it: enough for a decoder that
// believes the header to allocate all of it.
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

// TestDecodePictureRefusesADecompressionBomb is seventy bytes that say they
// are sixty thousand pixels on a side, which the standard decoder allocates
// fourteen gigabytes for before it finds there is nothing in them.
func TestDecodePictureRefusesADecompressionBomb(t *testing.T) {
	_, err := decodePicture(nil, "bomb.png", pngClaiming(60000, 60000))
	if err == nil || !strings.Contains(err.Error(), "more than") {
		t.Fatalf("the bomb decoded with %v, want it refused for its size", err)
	}
}

// TestDecodePictureRefusesAGIFOfTooManyFrames is the bomb spread over time:
// each frame is small in the file and a whole screen once it is read, so a GIF
// of a thousand frames on a big screen is refused for all of them together.
func TestDecodePictureRefusesAGIFOfTooManyFrames(t *testing.T) {
	pal := color.Palette{color.Black, color.White}
	g := &gif.GIF{Config: image.Config{Width: 4000, Height: 4000, ColorModel: pal}}
	for range 20 {
		g.Image = append(g.Image, image.NewPaletted(image.Rect(0, 0, 1, 1), pal))
		g.Delay = append(g.Delay, 1)
	}
	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, g); err != nil {
		t.Fatal(err)
	}
	if _, err := decodePicture(nil, "bomb.gif", buf.Bytes()); err == nil {
		t.Fatal("twenty frames of sixteen million pixels decoded, want them refused")
	}
}

// TestGIFFrameCountCountsEveryFrame keeps the walk honest against what the
// standard encoder writes, local colour tables and extensions included.
func TestGIFFrameCountCountsEveryFrame(t *testing.T) {
	pal := color.Palette{color.Black, color.White}
	g := &gif.GIF{}
	for range 7 {
		g.Image = append(g.Image, image.NewPaletted(image.Rect(0, 0, 3, 3), pal))
		g.Delay = append(g.Delay, 5)
	}
	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, g); err != nil {
		t.Fatal(err)
	}
	if n, err := gifFrameCount(buf.Bytes()); err != nil || n != 7 {
		t.Fatalf("counted %d frames (%v), want 7", n, err)
	}
	if _, err := decodePicture(nil, "small.gif", buf.Bytes()); err != nil {
		t.Fatalf("a small animated GIF was refused: %v", err)
	}
}
