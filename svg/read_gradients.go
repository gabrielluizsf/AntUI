package svg

// readGradients walks the drawing for the gradients it declares, keyed by the id
// a `url(#id)` names, and then resolves each one against the rest: an href takes
// what it does not say from another gradient, and only the whole drawing can say
// what that other one had.
func (img *Image) readGradients(root *element) {
	raw := map[string]*gradient{}
	collectGradients(root, img, raw)
	for id, g := range raw {
		img.grads[id] = g.resolved(raw)
	}
}
