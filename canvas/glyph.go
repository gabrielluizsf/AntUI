package canvas

import "github.com/gabrielluizsf/antui/font"

func (f *Face) glyph(r rune) *font.Mask {
	if f.file.Cmap.Glyph(r) == 0 {
		if n := f.fallbackFor(r); n != nil {
			return n.glyph(r)
		}
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if m, ok := f.glyphs[r]; ok {
		return m
	}
	g := f.file.Cmap.Glyph(r)
	m := font.Render(f.file.GlyphContours(g, 0), f.scale)
	f.glyphs[r] = &m
	return &m
}
