package svg

// maskRegion is where a `<mask>` reaches: a rectangle, either in fractions of
// the box of the element it is put on or in the drawing's own coordinates,
// which is what `maskUnits` chooses between. What is written against neither
// is the SVG default of the box grown by a tenth all round, filled in when the
// mask is read. What the mask says beyond its region is beyond the region
// rather than under it, so it keeps nothing there — the region is where the
// mask is allowed to say anything at all.
type maskRegion struct {
	// boxUnits is whether x, y, w and h are fractions of the box of the
	// element the mask is put on rather than lengths of the drawing's own.
	boxUnits   bool
	x, y, w, h float64
}
