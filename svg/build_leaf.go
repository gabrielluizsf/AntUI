package svg

// buildLeaf builds the two elements that are not an outline with a paint over
// them — writing and a picture — and answers nil for everything else. Writing
// is a piece of writing with a place and a size, kept as the pieces it was
// written in rather than as a shape; a picture is a bitmap, fetched while the
// drawing is read and fitted into its box when it is painted. Neither has
// anything inside it to build nodes for.
func (img *Image) buildLeaf(n *Node, e *element, st Style, warn func(string, ...any),
	def *maskDef) *Node {
	switch e.Name {
	case "text", "tspan":
		leaf := img.textNode(e, st, warn)
		leaf.clip = clipUnder(n.clip, st, leaf, warn)
		leaf.filters = n.filters
		leaf.mask = maskUnder(def, st, leaf, warn)
		return leaf
	case "image":
		leaf := img.imageNode(e, st, warn)
		leaf.clip = clipUnder(n.clip, st, leaf, warn)
		leaf.filters = n.filters
		leaf.mask = maskUnder(def, st, leaf, warn)
		return leaf
	}
	return nil
}
