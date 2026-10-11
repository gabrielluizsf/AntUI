package canvas

import "github.com/gabrielluizsf/antui/font"

// colorScale is how much a picture cut at its strike's size is enlarged by to
// draw at this face's own size.
func (f *Face) colorScale(cg *font.ColorGlyph) float64 {
	if cg == nil || cg.PPEM <= 0 {
		return 1
	}
	return f.size / float64(cg.PPEM)
}
