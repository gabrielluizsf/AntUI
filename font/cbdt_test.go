package font

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/assert"
)

// These tests write the CBLC/CBDT pair themselves: a real TrueType file with
// two tables put in its place, one strike at one size, and a PNG behind every
// glyph that strike names. The pair is read back through ParseTTF, so the
// directory and the checksums have to agree as well as the tables do.

// colorGlyphSpec is one picture in the test font's CBDT table, in the order
// the table keeps them.
type colorGlyphSpec struct {
	glyph int
	png   []byte
	// The five bytes the formats that carry metrics keep beside the
	// picture: the height and width it is cut at, the left and top
	// bearings, and how far the pen moves on — all in the pixels of the
	// strike it belongs to.
	height, width               int
	bearingX, bearingY, advance int
}

// cblcStrike is one index subtable in the test font's CBLC table.
type cblcStrike struct {
	ppem        int
	indexFormat int
	imageFormat int
	first, last int
	// Formats 1 and 3 keep one offset per glyph of the range and one past
	// the end of it, relative to the subtable's imageDataOffset. Format 4
	// pairs the offsets with the ids in ids in place of taking them by
	// position.
	offsets []int
	ids     []int
	// Formats 2 and 5 hold one size and one set of metrics for the whole
	// range, which the pictures of image format 19 have none of.
	imageSize int
	metrics   []byte
}

func put16(b []byte, at, v int) { b[at], b[at+1] = byte(v>>8), byte(v) }

func put32(b []byte, at, v int) {
	b[at], b[at+1], b[at+2], b[at+3] = byte(v>>24), byte(v>>16), byte(v>>8), byte(v)
}

// pngBits is one solid picture of the given size. The colour is what the
// tests read back out of the decoded image, to be sure they were handed the
// picture they asked for.
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

// bigMetrics is the eight bytes of a BigGlyphMetrics record: the horizontal
// half every run of text reads, then three for the vertical one.
func bigMetrics(s colorGlyphSpec) []byte {
	return []byte{
		byte(s.height), byte(s.width), byte(s.bearingX), byte(s.bearingY), byte(s.advance),
		0, 0, 0,
	}
}

// buildCBDT writes the picture data: the table's version, then one entry per
// glyph in the order given. The offsets it returns say where every entry
// begins relative to the imageDataOffset the subtables name, and one past the
// last, which is the list an index of format 1 or 3 keeps.
func buildCBDT(specs []colorGlyphSpec, imageFormat int) ([]byte, []int) {
	data := []byte{0, 3, 0, 0} // version 3.0
	offsets := []int{0}
	for _, s := range specs {
		switch imageFormat {
		case 17: // small metrics: the horizontal five bytes
			data = append(data, byte(s.height), byte(s.width),
				byte(s.bearingX), byte(s.bearingY), byte(s.advance))
		case 18: // big metrics: those five, then three for the vertical half
			data = append(data, bigMetrics(s)...)
		case 19: // no metrics: the index keeps them
		}
		n := len(s.png)
		data = append(data, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
		data = append(data, s.png...)
		offsets = append(offsets, len(data)-4)
	}
	return data, offsets
}

// buildIndexSubtable writes the header and the body of one index format. The
// header's imageDataOffset is where in the CBDT table the data begins: these
// tests keep it right behind the table's own version, so the offsets are
// counted from there.
func buildIndexSubtable(s cblcStrike) []byte {
	putHeader := func(b []byte) {
		put16(b, 0, s.indexFormat)
		put16(b, 2, s.imageFormat)
		put32(b, 4, 4)
	}
	switch s.indexFormat {
	case 1, 3:
		width := 4
		if s.indexFormat == 3 {
			width = 2
		}
		out := make([]byte, 8+width*len(s.offsets))
		putHeader(out)
		for i, o := range s.offsets {
			if width == 4 {
				put32(out, 8+i*width, o)
			} else {
				put16(out, 8+i*width, o)
			}
		}
		return out
	case 2:
		out := make([]byte, 20)
		putHeader(out)
		put32(out, 8, s.imageSize)
		copy(out[12:20], s.metrics)
		return out
	case 4:
		// One record per glyph and one past them all, pairing every id
		// with the offset where the glyph's data begins.
		out := make([]byte, 12+(len(s.ids)+1)*4)
		putHeader(out)
		put32(out, 8, len(s.ids))
		for i, g := range s.ids {
			put16(out, 12+i*4, g)
			put16(out, 14+i*4, s.offsets[i])
		}
		put16(out, 12+len(s.ids)*4, 0xFFFF)
		put16(out, 14+len(s.ids)*4, s.offsets[len(s.ids)])
		return out
	case 5:
		// The size and the metrics the pictures share, the count, then
		// the ids they belong to.
		out := make([]byte, 24+len(s.ids)*2)
		putHeader(out)
		put32(out, 8, s.imageSize)
		copy(out[12:20], s.metrics)
		put32(out, 20, len(s.ids))
		for i, g := range s.ids {
			put16(out, 24+i*2, g)
		}
		return out
	}
	return nil
}

// buildCBLC writes one BitmapSize record per strike, then every strike's
// list of index subtables behind them.
func buildCBLC(strikes []cblcStrike) []byte {
	subs := make([][]byte, len(strikes))
	for i, s := range strikes {
		subs[i] = buildIndexSubtable(s)
	}
	// The records come first at a fixed stride, so every list's place is
	// known before a byte of them is written.
	listOff := make([]int, len(strikes))
	off := 8 + 48*len(strikes)
	for i, sub := range subs {
		listOff[i] = off
		off += 8 + len(sub) // the record the subtable sits behind
	}
	out := make([]byte, off)
	put16(out, 0, 3) // version 3.0
	put16(out, 2, 0)
	put32(out, 4, len(strikes))
	for i, s := range strikes {
		rec := 8 + i*48
		put32(out, rec+0, listOff[i])
		put32(out, rec+4, 8+len(subs[i])) // the list: its record and subtable
		put32(out, rec+8, 1)              // index subtables the list holds
		put32(out, rec+12, 0)             // colorRef, which nothing reads
		// The two line metrics at rec+16 and rec+28 describe a run of
		// text across and down the strike; nothing here reads them.
		put16(out, rec+40, s.first)
		put16(out, rec+42, s.last)
		out[rec+44] = byte(s.ppem)
		out[rec+45] = byte(s.ppem)
		out[rec+46] = 32 // bit depth
		out[rec+47] = 1  // HORIZONTAL_METRICS
		put16(out, listOff[i], s.first)
		put16(out, listOff[i]+2, s.last)
		put32(out, listOff[i]+4, 8) // the subtable behind this record
		copy(out[listOff[i]+8:], subs[i])
	}
	return out
}

// checksum is the sum of a table's words, which the directory keeps beside
// every table it names. A table whose length is not a whole number of words
// is padded with zeroes, which add nothing.
func checksum(data []byte) int {
	var sum int
	for i := 0; i < len(data); i += 4 {
		var w [4]byte
		copy(w[:], data[i:min(i+4, len(data))])
		sum += int(be32(w[:], 0))
	}
	return sum
}

// withTables writes the named tables into a font file, taking the place of
// any table already of that name, and makes the directory agree with where
// everything landed.
func withTables(t *testing.T, base []byte, extra map[string][]byte) []byte {
	t.Helper()
	type table struct {
		name string
		data []byte
	}
	var tables []table
	have := map[string]bool{}
	for i := range int(be16(base, 4)) {
		rec := 12 + i*16
		name := string(base[rec : rec+4])
		d := base[int(be32(base, rec+8)) : int(be32(base, rec+8))+int(be32(base, rec+12))]
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
	// searchRange, entrySelector and rangeShift count the directory in
	// powers of two: the largest that fits, and what is left over.
	p := 0
	for 1<<(p+1) <= len(tables) {
		p++
	}
	put16(out, 6, 16<<p)
	put16(out, 8, p)
	put16(out, 10, 16*len(tables)-16<<p)
	for i, tb := range tables {
		rec := 12 + i*16
		copy(out[rec:], tb.name)
		put32(out, rec+4, checksum(tb.data))
		put32(out, rec+8, len(out))
		put32(out, rec+12, len(tb.data))
		out = append(out, tb.data...)
		for len(out)%4 != 0 {
			out = append(out, 0)
		}
	}
	return out
}

func roboto(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "testdata", "fonts", "Roboto-Regular.ttf"))
	if err != nil {
		t.Fatalf("read the test font: %v", err)
	}
	return b
}

// colorFont is a real TrueType file with a CBLC/CBDT pair written over it.
func colorFont(t *testing.T, strikes []cblcStrike, cbdt []byte) *TTF {
	t.Helper()
	with := withTables(t, roboto(t), map[string][]byte{
		"CBLC": buildCBLC(strikes),
		"CBDT": cbdt,
	})
	f, err := ParseTTF(with)
	if err != nil {
		t.Fatalf("ParseTTF: %v", err)
	}
	return f
}

func TestColorGlyphReadsAPictureAndItsMetrics(t *testing.T) {
	red := pngBits(t, 4, 4, color.NRGBA{R: 255, A: 255})
	blue := pngBits(t, 6, 3, color.NRGBA{B: 255, A: 255})
	specs := []colorGlyphSpec{
		{glyph: 7, png: red, height: 4, width: 4, bearingX: 1, bearingY: 5, advance: 6},
		{glyph: 8, png: blue, height: 3, width: 6, bearingX: -2, bearingY: 4, advance: 7},
	}
	cbdt, offsets := buildCBDT(specs, 17)
	f := colorFont(t, []cblcStrike{{
		ppem: 16, indexFormat: 1, imageFormat: 17, first: 7, last: 8, offsets: offsets,
	}}, cbdt)

	assert.True(t, f.HasColorGlyphs())

	g := f.ColorGlyph(7, 16)
	if g == nil {
		t.Fatal("want a picture for glyph 7")
	}
	assert.Equal(t, 16, g.PPEM)
	assert.Equal(t, 1, g.BearingX)
	assert.Equal(t, 5, g.BearingY)
	assert.Equal(t, 6, g.Advance)
	assert.Equal(t, 4, g.Image.Bounds().Dx())
	assert.Equal(t, 4, g.Image.Bounds().Dy())
	r, _, _, a := g.Image.At(0, 0).RGBA()
	assert.Equal(t, uint32(0xFFFF), r)
	assert.Equal(t, uint32(0xFFFF), a)

	// The second picture begins where the first one ends, which is the
	// only thing the offsets between them say: a wrong one hands back
	// the wrong PNG, or none.
	g2 := f.ColorGlyph(8, 16)
	if g2 == nil {
		t.Fatal("want a picture for glyph 8")
	}
	assert.Equal(t, -2, g2.BearingX)
	assert.Equal(t, 4, g2.BearingY)
	assert.Equal(t, 7, g2.Advance)
	assert.Equal(t, 6, g2.Image.Bounds().Dx())
	assert.Equal(t, 3, g2.Image.Bounds().Dy())
	_, _, b, _ := g2.Image.At(0, 0).RGBA()
	assert.Equal(t, uint32(0xFFFF), b)

	// A glyph the strike holds no picture for has none.
	assert.True(t, f.ColorGlyph(9, 16) == nil)
	// Nor does one no strike holds a size for beyond what the file cuts.
	assert.True(t, f.ColorGlyph(-1, 16) == nil)
}

func TestColorGlyphTakesTheMetricsFromTheIndex(t *testing.T) {
	specs := []colorGlyphSpec{{glyph: 11, png: pngBits(t, 5, 5, color.NRGBA{G: 255, A: 255}),
		height: 5, width: 5, bearingX: 1, bearingY: 5, advance: 8}}
	cbdt, _ := buildCBDT(specs, 19) // format 19 keeps no metrics of its own

	for _, tc := range []struct {
		name        string
		indexFormat int
	}{{"format 2", 2}, {"format 5", 5}} {
		t.Run(tc.name, func(t *testing.T) {
			f := colorFont(t, []cblcStrike{{
				ppem: 32, indexFormat: tc.indexFormat, imageFormat: 19,
				first: 11, last: 11, ids: []int{11},
				imageSize: len(cbdt) - 4, metrics: bigMetrics(specs[0]),
			}}, cbdt)
			g := f.ColorGlyph(11, 32)
			if g == nil {
				t.Fatal("want a picture for glyph 11")
			}
			assert.Equal(t, 32, g.PPEM)
			assert.Equal(t, 1, g.BearingX)
			assert.Equal(t, 5, g.BearingY)
			assert.Equal(t, 8, g.Advance)
			assert.Equal(t, 5, g.Image.Bounds().Dx())
		})
	}
}

func TestColorGlyphReadsEveryIndexFormat(t *testing.T) {
	specs := []colorGlyphSpec{{glyph: 5, png: pngBits(t, 4, 4, color.NRGBA{R: 255, G: 128, A: 255}),
		height: 4, width: 4, bearingX: 1, bearingY: 4, advance: 5}}
	for _, tc := range []struct {
		name                     string
		indexFormat, imageFormat int
	}{
		{"one offset per glyph", 1, 17},
		{"one size for the whole range", 2, 19},
		{"offsets of sixteen bits", 3, 17},
		{"glyphs named with their offsets", 4, 18},
		{"glyphs named by their ids", 5, 19},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cbdt, offsets := buildCBDT(specs, tc.imageFormat)
			strike := cblcStrike{ppem: 16, indexFormat: tc.indexFormat,
				imageFormat: tc.imageFormat, first: 5, last: 5, offsets: offsets}
			if tc.indexFormat >= 4 {
				strike.ids = []int{5}
			}
			if tc.indexFormat == 2 || tc.indexFormat == 5 {
				strike.imageSize = len(cbdt) - 4
				strike.metrics = bigMetrics(specs[0])
			}
			g := colorFont(t, []cblcStrike{strike}, cbdt).ColorGlyph(5, 16)
			if g == nil {
				t.Fatal("want a picture for glyph 5")
			}
			assert.Equal(t, 1, g.BearingX)
			assert.Equal(t, 4, g.BearingY)
			assert.Equal(t, 5, g.Advance)
			assert.Equal(t, 4, g.Image.Bounds().Dx())
		})
	}
}

func TestColorGlyphTakesTheStrikeNearestTheSize(t *testing.T) {
	// A strike's range takes the pictures one after the other, so the
	// small glyph 3 and glyph 4 sit in front of the big one.
	specs := []colorGlyphSpec{
		{glyph: 3, png: pngBits(t, 4, 4, color.NRGBA{R: 255, A: 255}),
			height: 4, width: 4, bearingX: 1, bearingY: 5, advance: 6},
		{glyph: 4, png: pngBits(t, 4, 4, color.NRGBA{G: 255, A: 255}),
			height: 4, width: 4, bearingX: 3, bearingY: 5, advance: 6},
		{glyph: 3, png: pngBits(t, 8, 8, color.NRGBA{B: 255, A: 255}),
			height: 8, width: 8, bearingX: 2, bearingY: 9, advance: 10},
	}
	cbdt, offsets := buildCBDT(specs, 17)
	f := colorFont(t, []cblcStrike{
		{ppem: 16, indexFormat: 1, imageFormat: 17, first: 3, last: 4,
			offsets: offsets[:3]},
		{ppem: 32, indexFormat: 1, imageFormat: 17, first: 3, last: 3,
			offsets: offsets[2:]},
	}, cbdt)

	// The strike cut closest to the size asked for is the one drawn.
	assert.Equal(t, 16, f.ColorGlyph(3, 16).PPEM)
	assert.Equal(t, 32, f.ColorGlyph(3, 32).PPEM)
	assert.Equal(t, 16, f.ColorGlyph(3, 20).PPEM)
	assert.Equal(t, 1, f.ColorGlyph(3, 16).BearingX)
	assert.Equal(t, 2, f.ColorGlyph(3, 32).BearingX)

	// A size the file did not cut takes the one it has, rather than no
	// picture at all.
	g := f.ColorGlyph(4, 64)
	if g == nil {
		t.Fatal("want the picture glyph 4 was cut at")
	}
	assert.Equal(t, 16, g.PPEM)
}

func TestColorGlyphIsAbsentWhenTheTableWillNotRead(t *testing.T) {
	t.Run("a list cut short", func(t *testing.T) {
		cbdt, offsets := buildCBDT([]colorGlyphSpec{{glyph: 5,
			png: pngBits(t, 4, 4, color.NRGBA{A: 255})}}, 17)
		cblc := buildCBLC([]cblcStrike{{ppem: 16, indexFormat: 1, imageFormat: 17,
			first: 5, last: 5, offsets: offsets}})
		f, err := ParseTTF(withTables(t, roboto(t), map[string][]byte{
			"CBLC": cblc[:len(cblc)/2],
			"CBDT": cbdt,
		}))
		if err != nil {
			t.Fatalf("ParseTTF: %v", err)
		}
		assert.False(t, f.HasColorGlyphs())
	})

	t.Run("an index format nothing draws", func(t *testing.T) {
		f := colorFont(t, []cblcStrike{{ppem: 16, indexFormat: 0, imageFormat: 17,
			first: 5, last: 5, offsets: []int{0, 1}}}, []byte{0, 3, 0, 0})
		assert.False(t, f.HasColorGlyphs())
	})

	t.Run("a picture of no colour", func(t *testing.T) {
		// Image format 1 is the monochrome one: a mask, not a PNG.
		f := colorFont(t, []cblcStrike{{ppem: 16, indexFormat: 1, imageFormat: 1,
			first: 5, last: 5, offsets: []int{0, 4}}}, []byte{0, 3, 0, 0, 1, 2, 3, 4})
		assert.False(t, f.HasColorGlyphs())
	})

	t.Run("a table that is not there", func(t *testing.T) {
		f, err := ParseTTF(roboto(t))
		if err != nil {
			t.Fatalf("ParseTTF: %v", err)
		}
		assert.False(t, f.HasColorGlyphs())
		assert.True(t, f.ColorGlyph(f.Cmap.Glyph('A'), 16) == nil)
	})

	t.Run("a picture with no metrics beside it", func(t *testing.T) {
		// Image format 19 carries no metrics of its own: it is written
		// against an index of format 2 or 5, which keeps them.
		specs := []colorGlyphSpec{{glyph: 5, png: pngBits(t, 4, 4, color.NRGBA{R: 255, A: 255})}}
		cbdt, offsets := buildCBDT(specs, 19)
		f := colorFont(t, []cblcStrike{{ppem: 16, indexFormat: 1, imageFormat: 19,
			first: 5, last: 5, offsets: offsets}}, cbdt)
		assert.True(t, f.HasColorGlyphs())
		assert.True(t, f.ColorGlyph(5, 16) == nil)
	})

	t.Run("a picture that is not a PNG", func(t *testing.T) {
		cbdt, offsets := buildCBDT([]colorGlyphSpec{{glyph: 5,
			png: []byte("not a png at all")}}, 17)
		f := colorFont(t, []cblcStrike{{ppem: 16, indexFormat: 1, imageFormat: 17,
			first: 5, last: 5, offsets: offsets}}, cbdt)
		assert.True(t, f.HasColorGlyphs())
		assert.True(t, f.ColorGlyph(5, 16) == nil)
	})
}
