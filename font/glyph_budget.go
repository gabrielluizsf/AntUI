package font

// glyphBudget is what one glyph may still read, counted through every
// composite inside it: the pieces it is put together from and the points
// those pieces bring. The nesting is bounded already, but not how wide each
// level is: a glyph of twenty copies of itself is 130 bytes that take
// seconds, and forty take minutes — the instancing bomb of an SVG `<use>`,
// played through a font. A thousand copies of a piece of a hundred thousand
// points is the same trick spent on memory instead.
type glyphBudget struct {
	components, points int
}

// An accented letter is two or three pieces and the busiest CJK composites a
// few dozen, and the most detailed outlines in real fonts a few thousand
// points; the rest of a glyph past these is a font built to stall whatever
// draws it.
const (
	maxComponents  = 1024
	maxGlyphPoints = 1 << 18
)

// spendPoints takes the points of the contours read from what the glyph may
// still read, and says whether there were that many left.
func (b *glyphBudget) spendPoints(contours []Contour) bool {
	for _, c := range contours {
		b.points -= len(c.X)
	}
	return b.points >= 0
}
