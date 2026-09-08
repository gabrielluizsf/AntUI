package font

import (
	"encoding/binary"
)

type CmapTable struct {
	Format uint16
	Data   []byte

	// format 4
	Segments int
	// format 12
	Groups int
}

// Glyph is the glyph index for a codepoint, or 0 — which is the font's own
// "not there" glyph, usually an empty box, and is what should be drawn.
func (t CmapTable) Glyph(r rune) int {
	switch t.Format {
	case 4:
		return t.glyph4(r)
	case 12:
		return t.glyph12(r)
	case 6:
		return t.glyph6(r)
	case 0:
		return t.glyph0(r)
	}
	return 0
}

func (t CmapTable) glyph4(r rune) int {
	if r > 0xFFFF || t.Segments == 0 {
		return 0
	}
	code := uint16(r)
	const endsAt = 14
	starts := endsAt + t.Segments*2 + 2
	deltas := starts + t.Segments*2
	ranges := deltas + t.Segments*2

	// A linear walk. The format promises the segments are sorted and a
	// binary search is the usual answer, but a face caches every glyph it
	// draws and there are rarely more than a hundred segments — so this runs
	// once per distinct character in the program's whole life.
	for i := range t.Segments {
		end := be16(t.Data, endsAt+i*2)
		if code > end {
			continue
		}
		start := be16(t.Data, starts+i*2)
		if code < start {
			return 0
		}
		delta := be16(t.Data, deltas+i*2)
		offset := be16(t.Data, ranges+i*2)
		if offset == 0 {
			return int(uint16(code + delta))
		}
		// The offset is from the position of the offset itself, which is the
		// one genuinely strange thing about this format.
		at := ranges + i*2 + int(offset) + int(code-start)*2
		if at+2 > len(t.Data) {
			return 0
		}
		g := be16(t.Data, at)
		if g == 0 {
			return 0
		}
		return int(uint16(g + delta))
	}
	return 0
}

func (t CmapTable) glyph12(r rune) int {
	for i := range t.Groups {
		g := 16 + i*12
		if g+12 > len(t.Data) {
			return 0
		}
		start, end := be32(t.Data, g), be32(t.Data, g+4)
		if uint32(r) < start {
			return 0
		}
		if uint32(r) > end {
			continue
		}
		return int(be32(t.Data, g+8) + uint32(r) - start)
	}
	return 0
}

func (t CmapTable) glyph6(r rune) int {
	if len(t.Data) < 10 {
		return 0
	}
	first, count := int(be16(t.Data, 6)), int(be16(t.Data, 8))
	i := int(r) - first
	if i < 0 || i >= count || 10+i*2+2 > len(t.Data) {
		return 0
	}
	return int(be16(t.Data, 10+i*2))
}

func be16(b []byte, i int) uint16 {
	if i+2 > len(b) {
		return 0
	}
	return binary.BigEndian.Uint16(b[i:])
}

func be32(b []byte, i int) uint32 {
	if i+4 > len(b) {
		return 0
	}
	return binary.BigEndian.Uint32(b[i:])
}

func (t CmapTable) glyph0(r rune) int {
	if r < 0 || r > 255 || len(t.Data) < 6+256 {
		return 0
	}
	return int(t.Data[6+int(r)])
}
