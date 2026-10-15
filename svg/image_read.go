package svg

// read turns the element tree into the node tree, carrying the style down as it
// goes so that a node's fill is the one it was given or, failing that, the one
// it inherited from the element above.
func (img *Image) read(e *element) {
	img.Root.Kids = append(img.Root.Kids, img.node(e, img.rootStyle()))
}
