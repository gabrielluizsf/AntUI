package canvas

import "github.com/gabrielluizsf/antui/font"

// WithFallback is this face with next handed every rune this one has no shape
// for, the way a font-family list hands a character down its names to the face
// that has it. The result draws differently from what this one draws, so it
// cannot share the glyph cache; it does share the font file, and the caches of
// sizes at which either of them is met.
func (f *Face) WithFallback(next *Face) *Face {
	if f == nil {
		return next
	}
	if next == nil {
		return f
	}
	return &Face{
		unit:     f.unit,
		file:     f.file,
		size:     f.size,
		scale:    f.scale,
		ascent:   f.ascent,
		descent:  f.descent,
		height:   f.height,
		hasColor: f.hasColor,
		next:     next,
		glyphs:   map[rune]*font.Mask{},
	}
}
