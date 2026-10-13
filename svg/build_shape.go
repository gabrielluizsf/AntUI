package svg

// buildShape stores the geometry in the drawing's coordinates and prepares
// the values derived from its style for the first paint.
func (img *Image) buildShape(n *Node, e *element, st Style) {
	n.Path = img.shape(e, st)
	if n.Path == nil {
		return
	}
	n.sourcePath = clonePath(n.Path)
	if minX, minY, maxX, maxY, ok := n.sourcePath.Bounds(); ok {
		n.box = [4]float64{minX, minY, maxX - minX, maxY - minY}
	}
	n.rawDash = cloneDash(n.Style.Dash)
	refreshShape(n)
}
