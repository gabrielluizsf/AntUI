package font

// Colour glyphs: the CBDT/CBLC pair of tables, which hold a picture per glyph
// instead of an outline — what an emoji font is made of. CBLC says which
// sizes the file keeps bitmaps at and where each glyph's data begins, and
// CBDT holds the data itself: a PNG, in one of the three image formats that
// carry one (17 and 18 with the metrics beside the picture, 19 with the
// metrics in the index instead). A table that does not read leaves the font
// with no colour bitmaps rather than failing it — the outlines the same file
// carries still draw.

import (
	"bytes"
	"image"
	"image/png"
	"math"
)

// ColorGlyph is one glyph drawn as a colour picture. Every number on it is in
// the pixels of the strike the picture was cut at: the top-left corner sits
// BearingX to the right of the pen and BearingY above the baseline, Advance
// is how far the pen moves on, and PPEM is how many pixels per em that strike
// was — a face at another size scales all of them by its own size over PPEM.
type ColorGlyph struct {
	Image              image.Image
	PPEM               int
	BearingX, BearingY int
	Advance            int
}

// sbitMetrics is a bitmap's size and placement in the pixels of its strike.
// The horizontal bearings of a big metrics record are what a run of text
// uses, and they are its first five bytes — which is all a small record
// holds.
type sbitMetrics struct {
	Height             int
	Width              int
	BearingX, BearingY int
	Advance            int
}

// colorBitmap is where one glyph's image data begins in the CBDT table, and
// the metrics when they are the index's own: image format 19 keeps none of
// its own and is written against an index that has them.
type colorBitmap struct {
	off, length int
	imageFormat int
	metrics     sbitMetrics
	inIndex     bool
}

// colorStrike is one size the file keeps bitmaps at.
type colorStrike struct {
	ppem   int
	glyphs map[int]colorBitmap
}

// HasColorGlyphs reports whether the file carries colour bitmaps, which is
// what a face asks before it looks for one on any rune.
func (f *TTF) HasColorGlyphs() bool { return len(f.colorStrikes) > 0 }

// ColorGlyph is the colour picture for one glyph at a size in pixels per em,
// or nil when the font has none for it. The strike nearest the size is taken,
// and a glyph that strike does not hold is looked for in the others: a file
// need not cut every size for every glyph.
func (f *TTF) ColorGlyph(glyph int, pixels float64) *ColorGlyph {
	if len(f.colorStrikes) == 0 || glyph < 0 {
		return nil
	}
	best, bestDist := -1, math.MaxFloat64
	for i := range f.colorStrikes {
		s := &f.colorStrikes[i]
		if _, ok := s.glyphs[glyph]; !ok {
			continue
		}
		if d := math.Abs(float64(s.ppem) - pixels); d < bestDist {
			best, bestDist = i, d
		}
	}
	if best < 0 {
		return nil
	}
	s := &f.colorStrikes[best]
	b := s.glyphs[glyph]
	if b.off < 0 || b.off+b.length > len(f.cbdt) {
		return nil
	}
	m, data, ok := colorEntry(b, f.cbdt[b.off:b.off+b.length])
	if !ok {
		return nil
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	return &ColorGlyph{
		Image:    img,
		PPEM:     s.ppem,
		BearingX: m.BearingX,
		BearingY: m.BearingY,
		Advance:  m.Advance,
	}
}

// colorEntry reads one CBDT image: the metrics it carries, if any, and the
// PNG behind its dataLen. The three formats differ only in which metrics sit
// in front of the length.
func colorEntry(b colorBitmap, d []byte) (sbitMetrics, []byte, bool) {
	var m sbitMetrics
	var at int
	switch b.imageFormat {
	case 17: // small metrics, then the length, then the PNG
		m, at = hMetrics(d), 9
	case 18: // big metrics, then the length, then the PNG
		m, at = hMetrics(d), 12
	case 19: // the length and the PNG; the metrics are the index's
		m, at = b.metrics, 4
		if !b.inIndex {
			return m, nil, false
		}
	default:
		return m, nil, false
	}
	if at > len(d) {
		return m, nil, false
	}
	n := int(be32(d, at-4))
	if n < 0 || at+n > len(d) {
		return m, nil, false
	}
	return m, d[at : at+n], true
}

// hMetrics reads the horizontal half of a metrics record: the same five
// bytes whether the record is the small one or the hori half of the big one.
func hMetrics(d []byte) sbitMetrics {
	if len(d) < 5 {
		return sbitMetrics{}
	}
	return sbitMetrics{
		Height:   int(d[0]),
		Width:    int(d[1]),
		BearingX: int(int8(d[2])),
		BearingY: int(int8(d[3])),
		Advance:  int(d[4]),
	}
}

// parseColorStrikes reads a CBLC table and the CBDT table it points into,
// into one entry per strike per glyph. Anything that does not read — a
// truncated table, an index format with nothing this engine draws, an image
// format that carries no PNG — is left out rather than reported.
func parseColorStrikes(cblc, cbdt []byte) []colorStrike {
	if len(cblc) < 8 || len(cbdt) < 4 {
		return nil
	}
	numSizes := int(be32(cblc, 4))
	if numSizes <= 0 || numSizes > (len(cblc)-8)/48 {
		return nil
	}
	var strikes []colorStrike
	for i := range numSizes {
		if s, ok := parseColorStrike(cblc, cbdt, 8+i*48); ok {
			strikes = append(strikes, s)
		}
	}
	return strikes
}

// parseColorStrike reads one BitmapSize record and the index subtables it
// points at. A strike that holds no picture is not a strike this font has.
func parseColorStrike(cblc, cbdt []byte, rec int) (colorStrike, bool) {
	s := colorStrike{ppem: int(cblc[rec+44]), glyphs: map[int]colorBitmap{}}
	listOff, numSubs := int(be32(cblc, rec)), int(be32(cblc, rec+8))
	if s.ppem == 0 || listOff < 0 || numSubs <= 0 ||
		numSubs > (len(cblc)-listOff)/8 {
		return s, false
	}
	for i := range numSubs {
		// The record names its subtable from the start of the list, and
		// the range it covers comes off the record itself.
		at := listOff + i*8
		first, last := int(be16(cblc, at)), int(be16(cblc, at+2))
		sub := listOff + int(be32(cblc, at+4))
		for g, b := range strikeGlyphs(cblc, cbdt, sub, first, last) {
			s.glyphs[g] = b
		}
	}
	return s, len(s.glyphs) > 0
}

// strikeGlyphs reads one IndexSubTable: the glyphs its range covers and where
// each one's data begins. Only image formats that hold a PNG are read; the
// rest are the monochrome tables this engine draws no colour from.
func strikeGlyphs(cblc, cbdt []byte, sub, first, last int) map[int]colorBitmap {
	if sub < 0 || sub+8 > len(cblc) || first > last {
		return nil
	}
	indexFormat, imageFormat := int(be16(cblc, sub)), int(be16(cblc, sub+2))
	if !isPNGImage(imageFormat) {
		return nil
	}
	base := int(be32(cblc, sub+4))
	out := map[int]colorBitmap{}
	// add places one glyph's data, as long as it lies inside the table.
	add := func(g, off, end int) {
		if g < first || g > last || off < 0 || off >= end || end > len(cbdt) {
			return
		}
		out[g] = colorBitmap{off: off, length: end - off, imageFormat: imageFormat}
	}
	// constant is the one set of metrics every glyph of the range shares,
	// kept here rather than repeated in each glyph's data.
	var constant sbitMetrics
	constantMetrics := false

	switch indexFormat {
	case 1, 3:
		// One offset per glyph in the range, and one past the end of it
		// to measure the last: 32-bit offsets in format 1, 16-bit in
		// format 3, both added to the base in the header.
		width := 4
		if indexFormat == 3 {
			width = 2
		}
		num := last - first + 2
		if sub+8+num*width > len(cblc) {
			return nil
		}
		for i := 0; i < num-1; i++ {
			add(first+i,
				base+int(readWidth(cblc, sub+8+i*width, width)),
				base+int(readWidth(cblc, sub+8+(i+1)*width, width)))
		}
	case 2:
		// Every glyph the same size, all with the metrics held here.
		if sub+20 > len(cblc) {
			return nil
		}
		size := int(be32(cblc, sub+8))
		constant, constantMetrics = hMetrics(cblc[sub+12:]), true
		for g := first; g <= last; g++ {
			off := base + (g-first)*size
			add(g, off, off+size)
		}
	case 4:
		// Sparse glyphs, each named with its offset, and one record
		// past the last of them to measure it.
		if sub+12 > len(cblc) {
			return nil
		}
		count := int(be32(cblc, sub+8))
		if count <= 0 || sub+12+(count+1)*4 > len(cblc) {
			return nil
		}
		for i := 0; i < count; i++ {
			at := sub + 12 + i*4
			// The record after this one is where this glyph's data
			// ends: its offset, not the id beside it.
			next := sub + 12 + (i+1)*4
			add(int(be16(cblc, at)),
				base+int(be16(cblc, at+2)),
				base+int(be16(cblc, next+2)))
		}
	case 5:
		// Sparse glyphs with one size and one set of metrics for the
		// whole range: imageSize, the metrics, the count, then the ids.
		if sub+24 > len(cblc) {
			return nil
		}
		size := int(be32(cblc, sub+8))
		count := int(be32(cblc, sub+20))
		if count <= 0 || sub+24+count*2 > len(cblc) {
			return nil
		}
		constant, constantMetrics = hMetrics(cblc[sub+12:]), true
		for i := range count {
			off := base + i*size
			add(int(be16(cblc, sub+24+i*2)), off, off+size)
		}
	default:
		return nil
	}
	if constantMetrics {
		for g, b := range out {
			b.metrics, b.inIndex = constant, true
			out[g] = b
		}
	}
	return out
}

// isPNGImage reports whether a CBDT image format holds a PNG, which is the
// only colour picture this engine reads.
func isPNGImage(format int) bool { return format == 17 || format == 18 || format == 19 }

// readWidth reads one offset from an index subtable, 16 or 32 bits wide.
func readWidth(b []byte, i, width int) uint32 {
	if width == 2 {
		return uint32(be16(b, i))
	}
	return be32(b, i)
}
