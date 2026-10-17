package svg

// selNode is one element with the neighbours a selector walks: the element
// above it, the ones written beside it, and the elements inside it. A
// descendant or a sibling combinator is answered by walking these rather than
// by the shape of the drawing.
type selNode struct {
	el     *element
	parent *selNode
	prev   *selNode
	next   *selNode
	kids   []*selNode
}
