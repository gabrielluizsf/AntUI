package canvas

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"slices"
	"strings"
	"testing"
)

// A font whose pictures stand in for outlines — what an emoji font is made
// of. These tests write a CBLC/CBDT pair over a TrueType file, read it as a
// face, and check that the picture rather than an outline of the glyph is
// what the face draws, measures and hands down a chain.

func put16(b []byte, at, v int) { b[at], b[at+1] = byte(v>>8), byte(v) }

func put32(b []byte, at, v int) {
	b[at], b[at+1], b[at+2], b[at+3] = byte(v>>24), byte(v>>16), byte(v>>8), byte(v)
}

func get16(b []byte, at int) int { return int(b[at])<<8 | int(b[at+1]) }

func get32(b []byte, at int) int {
	return int(b[at])<<24 | int(b[at+1])<<16 | int(b[at+2])<<8 | int(b[at+3])
}

// pngBits is one solid picture of the given size. The colour is what the
// tests read back out of the canvas, to be sure the picture is what drew
// there rather than an outline in the colour the caller asked for.
func pngBits(t *testing.T, w, h int, c color.NRGBA) []byte {
	t.Helper()
	m := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			m.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, m); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

// colorTables writes a CBLC/CBDT pair holding one picture for one glyph at
// one size: a single strike, the one index subtable it names, and an entry of
// image format 17 — five bytes of metrics, a length, then the PNG.
func colorTables(glyph, ppem int, pic []byte) (cblc, cbdt []byte) {
	entry := []byte{5, 5, 1, 6, 7} // height, width, bearingX, bearingY, advance
	n := len(pic)
	entry = append(entry, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
	entry = append(entry, pic...)
	cbdt = append([]byte{0, 3, 0, 0}, entry...) // version 3.0

	// The CBLC header, the one BitmapSize record behind it, then the list
	// that record names at 56, then the subtable that list names at 64.
	cblc = make([]byte, 80)
	put16(cblc, 0, 3) // version 3.0
	put32(cblc, 4, 1) // strikes the table holds
	put32(cblc, 8, 56)
	put32(cblc, 12, 16)
	put32(cblc, 16, 1)
	put16(cblc, 48, glyph)
	put16(cblc, 50, glyph)
	cblc[52], cblc[53] = byte(ppem), byte(ppem)
	cblc[54] = 32 // bit depth
	cblc[55] = 1  // horizontal metrics
	put16(cblc, 56, glyph)
	put16(cblc, 58, glyph)
	put32(cblc, 60, 8)
	put16(cblc, 64, 1)  // index format 1: an offset per glyph
	put16(cblc, 66, 17) // image format 17: metrics, a length, a PNG
	put32(cblc, 68, 4)  // the data begins behind the CBDT version
	put32(cblc, 72, 0)
	put32(cblc, 76, len(entry))
	return cblc, cbdt
}

// tableChecksum is the sum of a table's words, which the directory keeps
// beside every table it names. A table not a whole number of words long is
// padded with zeroes, which add nothing.
func tableChecksum(data []byte) int {
	var sum int
	for i := 0; i < len(data); i += 4 {
		var w [4]byte
		copy(w[:], data[i:min(i+4, len(data))])
		sum += get32(w[:], 0)
	}
	return sum
}

// withColorTables writes the pair into a font file, taking the place of any
// table already of that name, and makes the directory agree with where
// everything landed.
func withColorTables(t *testing.T, base []byte, cblc, cbdt []byte) []byte {
	t.Helper()
	type table struct {
		name string
		data []byte
	}
	extra := map[string][]byte{"CBLC": cblc, "CBDT": cbdt}
	var tables []table
	have := map[string]bool{}
	for i := range get16(base, 4) {
		rec := 12 + i*16
		name := string(base[rec : rec+4])
		d := base[get32(base, rec+8) : get32(base, rec+8)+get32(base, rec+12)]
		if w, ok := extra[name]; ok {
			d, have[name] = w, true
		}
		tables = append(tables, table{name, d})
	}
	for name, d := range extra {
		if !have[name] {
			tables = append(tables, table{name, d})
		}
	}
	slices.SortFunc(tables, func(a, b table) int { return strings.Compare(a.name, b.name) })

	out := make([]byte, 12+16*len(tables))
	copy(out, base[:4]) // sfnt version
	put16(out, 4, len(tables))
	for i, tb := range tables {
		rec := 12 + i*16
		copy(out[rec:], tb.name)
		put32(out, rec+4, tableChecksum(tb.data))
		put32(out, rec+8, len(out))
		put32(out, rec+12, len(tb.data))
		out = append(out, tb.data...)
		for len(out)%4 != 0 {
			out = append(out, 0)
		}
	}
	return out
}

// colorFace is a TrueType file with a picture standing in for one of its
// glyphs: the file at path, drawn at pixels per em, whose glyph r is a solid
// picture at a strike of sixteen.
func colorFace(t *testing.T, path string, r rune, pixels float64) *Face {
	t.Helper()
	gid := loadFace(t, path, pixels).file.Cmap.Glyph(r)
	if gid == 0 {
		t.Fatalf("%q has no glyph for %q", path, r)
	}
	base, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	cblc, cbdt := colorTables(gid, 16, pngBits(t, 5, 5, color.NRGBA{B: 255, A: 255}))
	f, err := ParseFace(withColorTables(t, base, cblc, cbdt), pixels)
	if err != nil {
		t.Fatalf("ParseFace: %v", err)
	}
	return f
}

// TestFaceDrawsAColorPicture checks a face whose glyph is a picture: the pen
// moves on by the picture's own advance, the picture lands where its
// bearings say, the colour it was cut with is what draws rather than the
// colour the caller handed, and a face asked for twice the size draws the
// picture twice as large.
func TestFaceDrawsAColorPicture(t *testing.T) {
	const path = "../testdata/fonts/Roboto-Regular.ttf"
	want := RGBA(0, 0, 255, 255)
	face := colorFace(t, path, 'A', 16)
	if face.colorGlyph('A') == nil {
		t.Fatal("the face holds no picture for A")
	}

	// The picture carries its own advance, which is what the face measures
	// a string by rather than the width of an outline it does not draw.
	if got, adv := face.Width("A"), 7; got != adv {
		t.Errorf("A is %d wide, want the picture's %d", got, adv)
	}
	cv, err := NewCanvas(32, 32)
	if err != nil {
		t.Fatal(err)
	}
	// The picture is five pixels across, one to the right of the pen and
	// six above the baseline Draw is given.
	if got := face.Draw(cv, 10, 20, "A", White); got != 7 {
		t.Errorf("the pen moved on by %d, want the picture's 7", got)
	}
	for _, p := range [][2]int{{11, 14}, {13, 16}, {15, 18}} {
		if got := cv.At(p[0], p[1]); got != want {
			t.Errorf("At(%d,%d) = %v, want the picture's %v", p[0], p[1], got, want)
		}
	}
	// Nothing is drawn where the picture is not, and not in the colour
	// Draw was handed either.
	for _, p := range [][2]int{{10, 14}, {16, 14}, {11, 13}, {16, 18}, {15, 19}} {
		if got := cv.At(p[0], p[1]); got != 0 {
			t.Errorf("At(%d,%d) = %v, want the bare canvas", p[0], p[1], got)
		}
	}

	// The picture was cut at a strike of sixteen pixels per em: a face at
	// twice that draws it twice as large, still from the same file.
	big := face.Scaled(2)
	if big.colorGlyph('A') == nil {
		t.Fatal("the face at twice the size holds no picture")
	}
	wide, err := NewCanvas(48, 48)
	if err != nil {
		t.Fatal(err)
	}
	if got := big.Draw(wide, 10, 40, "A", White); got != 14 {
		t.Errorf("the pen moved on by %d, want twice the picture's advance", got)
	}
	for _, p := range [][2]int{{12, 28}, {16, 32}, {21, 37}} {
		if got := wide.At(p[0], p[1]); got != want {
			t.Errorf("At(%d,%d) = %v, want the enlarged picture's %v", p[0], p[1], got, want)
		}
	}
	for _, p := range [][2]int{{11, 28}, {22, 37}, {12, 27}, {22, 28}} {
		if got := wide.At(p[0], p[1]); got != 0 {
			t.Errorf("At(%d,%d) = %v, want the bare canvas", p[0], p[1], got)
		}
	}

	// A style doubles the picture for a weight there is no file for and
	// leans it for a slant, the same as it does an outline.
	styled, err := NewCanvas(48, 48)
	if err != nil {
		t.Fatal(err)
	}
	baseline := 20 + face.Ascent()
	styled.DrawStyled(10, 20, "A", White, TextStyle{Face: face, Bold: true})
	if got := styled.At(16, baseline-6); got != want {
		t.Errorf("beside the picture = %v, want the doubled stroke's %v", got, want)
	}

	slanted, err := NewCanvas(48, 48)
	if err != nil {
		t.Fatal(err)
	}
	slanted.DrawStyled(10, 20, "A", White, TextStyle{Face: face, Italic: true})
	if got := slanted.At(11, baseline-6); got != 0 {
		t.Errorf("the top row of the picture did not lean: At(11,%d) = %v", baseline-6, got)
	}
	if got := slanted.At(12, baseline-6); got != want {
		t.Errorf("the leaned top row = %v, want the picture's %v", got, want)
	}
	if got := slanted.At(11, baseline-5); got != want {
		t.Errorf("the row under it = %v, want the picture's %v", got, want)
	}
}

// TestFaceTakesAColorPictureDownTheChain checks which family a rune's
// picture comes from: a rune the face has no glyph for reaches the face
// behind it and takes its picture, and a rune this face does have is drawn
// from this file — the names after the first in a font-family are for the
// runes it lacks.
func TestFaceTakesAColorPictureDownTheChain(t *testing.T) {
	const robotoPath = "../testdata/fonts/Roboto-Regular.ttf"
	const monoPath = "../testdata/fonts/DejaVuSansMono.ttf"
	want := RGBA(0, 0, 255, 255)

	roboto := loadFace(t, robotoPath, 16)
	if roboto.file.Cmap.Glyph('→') != 0 {
		t.Fatal("Roboto has an arrow — this test needs a rune it does not")
	}
	pictured := colorFace(t, monoPath, '→', 16)

	chained := roboto.WithFallback(pictured)
	if chained.colorGlyph('→') == nil {
		t.Fatal("the arrow took no picture from the face behind")
	}
	if got, adv := chained.Width("→"), 7; got != adv {
		t.Errorf("arrow width = %d, want the picture's %d", got, adv)
	}
	cv, err := NewCanvas(32, 32)
	if err != nil {
		t.Fatal(err)
	}
	if got := chained.Draw(cv, 10, 20, "→", White); got != 7 {
		t.Errorf("the pen moved on by %d, want the picture's 7", got)
	}
	if got := cv.At(11, 14); got != want {
		t.Errorf("the arrow was drawn as %v, want the picture from the face behind", got)
	}

	// A rune the first face has a glyph for is drawn from it, even when
	// the face behind it keeps a picture of the same rune.
	first := colorFace(t, monoPath, 'A', 16)
	joined := roboto.WithFallback(first)
	if joined.colorGlyph('A') != nil {
		t.Error("the face behind was drawn for a rune this one has")
	}
	if got, want := joined.Width("A"), roboto.Width("A"); got != want {
		t.Errorf("A is %d wide, want the outline's %d", got, want)
	}
}
