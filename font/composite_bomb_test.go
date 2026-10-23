package font

import (
	"encoding/binary"
	"testing"
	"time"
)

// selfComposite is a font of one glyph made of n copies of itself, each piece
// at the origin: n to the sixth glyphs read in a file of a few bytes a piece.
func selfComposite(n int) *TTF {
	g := make([]byte, 10)
	binary.BigEndian.PutUint16(g, 0xFFFF) // -1 contours: a composite
	for i := range n {
		flags := uint16(0x0002) // the arguments are an offset, a byte each
		if i < n-1 {
			flags |= 0x0020 // more pieces follow
		}
		g = binary.BigEndian.AppendUint16(g, flags)
		g = binary.BigEndian.AppendUint16(g, 0)
		g = append(g, 0, 0)
	}
	loca := make([]byte, 8)
	binary.BigEndian.PutUint32(loca[4:], uint32(len(g)))
	return &TTF{NumGlyphs: 1, LongLoca: true, Loca: loca, Glyf: g, UnitsPerEm: 1000}
}

// TestGlyphContoursStopsACompositeBomb is a 250-byte glyph of forty copies of
// itself, which without a limit on the pieces read is forty to the sixth —
// minutes of work for one letter.
func TestGlyphContoursStopsACompositeBomb(t *testing.T) {
	f := selfComposite(40)
	start := time.Now()
	f.GlyphContours(0, 0)
	if took := time.Since(start); took > time.Second {
		t.Fatalf("one glyph took %v, want it cut short", took)
	}
}
