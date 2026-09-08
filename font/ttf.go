package font

import (
	"fmt"
	"math"
)

// ParseTTF reads the tables a face needs out of a font file.
func ParseTTF(data []byte) (*TTF, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("antUI: %d bytes is not a font", len(data))
	}
	switch tag := be32(data, 0); tag {
	case 0x00010000, 0x74727565: // 1.0, and "true"
	case 0x74746366: // "ttcf": a collection, whose first face is taken
		if len(data) < 16 {
			return nil, fmt.Errorf("antUI: a truncated font collection")
		}
		offset := be32(data, 12)
		if int(offset) >= len(data) {
			return nil, fmt.Errorf("antUI: a font collection pointing past its end")
		}
		return parseTTFAt(data, int(offset))
	case 0x4F54544F: // "OTTO"
		return nil, fmt.Errorf("antUI: this is a CFF font, and only TrueType " +
			"outlines are read here")
	default:
		return nil, fmt.Errorf("antUI: %#08x is not a font's tag", tag)
	}
	return parseTTFAt(data, 0)
}

func parseTTFAt(data []byte, base int) (*TTF, error) {
	if base+12 > len(data) {
		return nil, fmt.Errorf("antUI: a truncated font")
	}
	count := int(be16(data, base+4))
	tables := map[string][]byte{}
	for i := range count {
		rec := base + 12 + i*16
		if rec+16 > len(data) {
			return nil, fmt.Errorf("antUI: a truncated table directory")
		}
		name := string(data[rec : rec+4])
		off, length := int(be32(data, rec+8)), int(be32(data, rec+12))
		if off < 0 || length < 0 || off+length > len(data) {
			// A table that runs off the end is skipped rather than fatal:
			// some fonts pad the last one, and refusing the whole file for a
			// table nothing here reads would be absurd.
			continue
		}
		tables[name] = data[off : off+length]
	}

	f := &TTF{Data: data}
	head, ok := tables["head"]
	if !ok || len(head) < 54 {
		return nil, fmt.Errorf("antUI: the font has no head table")
	}
	f.UnitsPerEm = be16(head, 18)
	if f.UnitsPerEm == 0 {
		return nil, fmt.Errorf("antUI: the font says it has no units per em")
	}
	f.LongLoca = be16(head, 50) == 1

	maxp, ok := tables["maxp"]
	if !ok || len(maxp) < 6 {
		return nil, fmt.Errorf("antUI: the font has no maxp table")
	}
	f.NumGlyphs = int(be16(maxp, 4))

	hhea, ok := tables["hhea"]
	if !ok || len(hhea) < 36 {
		return nil, fmt.Errorf("antUI: the font has no hhea table")
	}
	f.Ascent = int16(be16(hhea, 4))
	f.Descent = int16(be16(hhea, 6))
	f.LineGap = int16(be16(hhea, 8))
	f.NumHMetrics = int(be16(hhea, 34))

	f.Loca, f.Glyf, f.Hmtx = tables["loca"], tables["glyf"], tables["hmtx"]
	if f.Glyf == nil || f.Loca == nil {
		return nil, fmt.Errorf("antUI: the font has no TrueType outlines; a " +
			"CFF font cannot be drawn here")
	}
	cm, ok := tables["cmap"]
	if !ok {
		return nil, fmt.Errorf("antUI: the font has no cmap table")
	}
	c, err := parseCmap(cm)
	if err != nil {
		return nil, err
	}
	f.Cmap = c
	return f, nil
}

// parseCmap picks the best subtable in a cmap and prepares it.
//
// "Best" is a Unicode one, and among those the widest: a format 12 table can
// say things a format 4 table cannot, and a font that has both means the
// same by them for the characters they share.
func parseCmap(data []byte) (CmapTable, error) {
	if len(data) < 4 {
		return CmapTable{}, fmt.Errorf("antUI: a truncated cmap")
	}
	count := int(be16(data, 2))
	best, bestScore := -1, -1
	for i := range count {
		rec := 4 + i*8
		if rec+8 > len(data) {
			break
		}
		platform, encoding := be16(data, rec), be16(data, rec+2)
		offset := int(be32(data, rec+4))
		if offset >= len(data) {
			continue
		}
		score := -1
		switch {
		case platform == 3 && encoding == 10: // Windows, full Unicode
			score = 4
		case platform == 0: // Unicode, any version
			score = 3
		case platform == 3 && encoding == 1: // Windows, BMP
			score = 2
		case platform == 1 && encoding == 0: // Macintosh Roman
			score = 1
		}
		if score > bestScore {
			best, bestScore = offset, score
		}
	}
	if best < 0 {
		return CmapTable{}, fmt.Errorf("antUI: the font has no Unicode cmap")
	}

	sub := data[best:]
	if len(sub) < 4 {
		return CmapTable{}, fmt.Errorf("antUI: a truncated cmap subtable")
	}
	t := CmapTable{Format: be16(sub, 0), Data: sub}
	switch t.Format {
	case 4:
		if len(sub) < 14 {
			return CmapTable{}, fmt.Errorf("antUI: a truncated format 4 cmap")
		}
		t.Segments = int(be16(sub, 6)) / 2
	case 12:
		if len(sub) < 16 {
			return CmapTable{}, fmt.Errorf("antUI: a truncated format 12 cmap")
		}
		t.Groups = int(be32(sub, 12))
	default:
		return CmapTable{}, fmt.Errorf("antUI: cmap format %d is not read here", t.Format)
	}
	return t, nil
}

// TTF is a parsed font file. One file may hold several faces at several
// sizes; this is the file, and [Face] is one size of it.
type TTF struct {
	Data []byte

	UnitsPerEm uint16
	LongLoca   bool
	NumGlyphs  int

	Ascent, Descent, LineGap int16
	NumHMetrics              int

	Loca, Glyf, Hmtx []byte
	Cmap             CmapTable
}

// Advance is how far the pen moves after a glyph, in font units.
func (f *TTF) Advance(glyph int) int16 {
	if f.NumHMetrics == 0 || len(f.Hmtx) < 4 {
		return int16(f.UnitsPerEm / 2)
	}
	i := min(glyph, f.NumHMetrics-1)
	if (i+1)*4 > len(f.Hmtx) {
		return int16(f.UnitsPerEm / 2)
	}
	return int16(be16(f.Hmtx, i*4))
}

// Outline is where one glyph's contours are in the glyf table.
func (f *TTF) Outline(glyph int) []byte {
	if glyph < 0 || glyph >= f.NumGlyphs {
		return nil
	}
	var start, end int
	if f.LongLoca {
		if (glyph+2)*4 > len(f.Loca) {
			return nil
		}
		start, end = int(be32(f.Loca, glyph*4)), int(be32(f.Loca, (glyph+1)*4))
	} else {
		if (glyph+2)*2 > len(f.Loca) {
			return nil
		}
		start, end = int(be16(f.Loca, glyph*2))*2, int(be16(f.Loca, (glyph+1)*2))*2
	}
	if start >= end || end > len(f.Glyf) {
		return nil // an empty glyph, which a space is
	}
	return f.Glyf[start:end]
}

// Contour is one closed loop of a glyph, in font units, with each point
// marked as on or off the curve.
type Contour struct {
	X, Y []float64
	On   []bool
}

// GlyphContours reads one glyph's outline, following composites.
func (f *TTF) GlyphContours(glyph, depth int) []Contour {
	if depth > 5 {
		return nil
	}
	data := f.Outline(glyph)
	if len(data) < 10 {
		return nil
	}
	if n := int16(be16(data, 0)); n >= 0 {
		return simpleContours(data, int(n))
	}
	return f.compositeContours(data, depth)
}

// compositeContours reads a glyph made of other glyphs — an accented letter
// is a letter and an accent, placed.
func (f *TTF) compositeContours(data []byte, depth int) []Contour {
	var out []Contour
	p := 10
	for {
		if p+4 > len(data) {
			return out
		}
		flags := be16(data, p)
		index := int(be16(data, p+2))
		p += 4

		var dx, dy float64
		if flags&0x0001 != 0 { // ARG_1_AND_2_ARE_WORDS
			if p+4 > len(data) {
				return out
			}
			dx, dy = float64(int16(be16(data, p))), float64(int16(be16(data, p+2)))
			p += 4
		} else {
			if p+2 > len(data) {
				return out
			}
			dx, dy = float64(int8(data[p])), float64(int8(data[p+1]))
			p += 2
		}
		if flags&0x0002 == 0 {
			// The arguments are point numbers to match up rather than an
			// offset. Almost nothing uses it and guessing would place the
			// piece wrongly, so it is placed at the origin instead.
			dx, dy = 0, 0
		}

		a, b, c, d := 1.0, 0.0, 0.0, 1.0
		switch {
		case flags&0x0008 != 0: // WE_HAVE_A_SCALE
			if p+2 > len(data) {
				return out
			}
			a = f2dot14(be16(data, p))
			d = a
			p += 2
		case flags&0x0040 != 0: // X_AND_Y_SCALE
			if p+4 > len(data) {
				return out
			}
			a, d = f2dot14(be16(data, p)), f2dot14(be16(data, p+2))
			p += 4
		case flags&0x0080 != 0: // TWO_BY_TWO
			if p+8 > len(data) {
				return out
			}
			a, b = f2dot14(be16(data, p)), f2dot14(be16(data, p+2))
			c, d = f2dot14(be16(data, p+4)), f2dot14(be16(data, p+6))
			p += 8
		}

		for _, part := range f.GlyphContours(index, depth+1) {
			moved := Contour{
				X:  make([]float64, len(part.X)),
				Y:  make([]float64, len(part.Y)),
				On: part.On,
			}
			for i := range part.X {
				moved.X[i] = a*part.X[i] + c*part.Y[i] + dx
				moved.Y[i] = b*part.X[i] + d*part.Y[i] + dy
			}
			out = append(out, moved)
		}
		if flags&0x0020 == 0 { // MORE_COMPONENTS
			return out
		}
	}
}

// f2dot14 is the fixed-point form composite transforms are written in.
func f2dot14(v uint16) float64 { return float64(int16(v)) / 16384 }

// Scaled is how many pixels one font unit is at a size in pixels per em.
func (f *TTF) Scaled(pixels float64) float64 {
	return pixels / float64(f.UnitsPerEm)
}

// LineHeight is ascent + descent + gap, in pixels, rounded away from zero so
// that two lines never overlap by a fraction.
func (f *TTF) LineHeight(pixels float64) int {
	s := f.Scaled(pixels)
	h := (float64(f.Ascent) - float64(f.Descent) + float64(f.LineGap)) * s
	return int(math.Ceil(h))
}
