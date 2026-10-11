package svg

// readIds takes every element that wrote an id, under the id it wrote, for a
// `<use>` to be followed to what it names. The first element to an id is the
// one it answers to, the same as the gradients: a drawing that writes the same
// id twice is broken rather than undecided, and guessing which one was meant
// would be worse than picking the first.
func (img *Image) readIds(root *element) {
	img.ids = map[string]*element{}
	collectIds(root, img.ids)
}
