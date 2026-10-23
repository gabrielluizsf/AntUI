package font

import (
	"bytes"
	"image/png"
	"math"
)

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
	// The size a PNG says it is is allocated before a byte of it is
	// decoded, so it is asked first: seventy bytes may say they are sixty
	// thousand pixels on a side.
	cfg, err := png.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > maxColorGlyphSide || cfg.Height > maxColorGlyphSide {
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
