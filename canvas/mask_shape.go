package canvas

// A MaskShape is one outline a mask is cut from: the path, and the rule that
// says which of its insides count — the same rule a fill takes, because a mask
// is a fill measured rather than painted.
type MaskShape struct {
	Path *Path
	Rule FillRule
	// And is the cuts the shape itself carries: what it keeps has to be under
	// the shape and under every one of these too, each entry being the union
	// of one clip of shapes — the clip written inside a `<clipPath>` cutting
	// the shape it was written on. What any entry covers is inside that entry,
	// and the shape is kept only where the shape and every entry agree. An
	// entry with no shapes in it keeps nothing anywhere, so a shape cut by a
	// reference the drawing could not follow keeps nothing, the same as an
	// empty `<clipPath>` keeps nothing. Nil is a shape with no cuts at all.
	And [][]MaskShape
}
