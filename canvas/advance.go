package canvas

import "github.com/gabrielluizsf/antui/font"

// advance is how far one rune moves the pen, in this face's pixels: its own
// advance when it holds the rune, and the advance of the face in the chain
// that does when this one does not.
func (f *Face) advance(r rune) float64 {
	if f == nil || f.file == nil {
		return float64(font.TextWidth(string(r))) * f.scaleOf()
	}
	if cg := f.colorGlyph(r); cg != nil {
		return float64(cg.Advance) * f.colorScale(cg)
	}
	g := f.file.Cmap.Glyph(r)
	if g == 0 {
		if n := f.fallbackFor(r); n != nil {
			return n.advance(r)
		}
	}
	return float64(f.file.Advance(g)) * f.scale
}
