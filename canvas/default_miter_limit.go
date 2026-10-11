package canvas

// DefaultMiterLimit is the corner SVG gives up on once a join's miter grows past
// this many stroke widths. Four keeps the spike of a right angle and drops it
// on anything sharper.
const DefaultMiterLimit = 4
