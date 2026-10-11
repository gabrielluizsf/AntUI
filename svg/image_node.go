package svg

// node builds one node and the nodes inside it.
func (img *Image) node(e *element, inherited Style) *Node {
	return img.build(e, inherited, false)
}
