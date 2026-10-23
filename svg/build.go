package svg

// build builds one node and the nodes inside it. forceKids says that what is
// inside this element is drawn after all, which is the one case of a
// definition being taken out of where it stands and put into the picture: a
// `<use>` naming a `<symbol>`. Everywhere else a definition keeps its contents
// to itself, because they are material for pointing at rather than part of the
// picture.
func (img *Image) build(e *element, inherited Style, forceKids bool) *Node {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	if !img.spendNode(warn) {
		return &Node{Name: e.Name, Style: Style{Hidden: true}}
	}
	st := inherited.with(e, warn)
	if e.Name == "use" {
		// A `<use>` follows its own clip itself, after its x and y have moved
		// it and before it builds what it points at — see [Image.useNode].
		return img.useNode(e, st, warn)
	}
	clip, def, filters := img.resolveRefs(&st, warn)
	n := &Node{Name: e.Name, Style: st, clip: clip, filters: filters}
	if n.Style.Hidden {
		return n
	}
	if leaf := img.buildLeaf(n, e, st, warn, def); leaf != nil {
		return leaf
	}
	img.buildShape(n, e, st)
	if !forceKids && isKnownButUnpainted(e.Name) {
		// What is written inside a `<defs>`, a `<clipPath>` or a `<symbol>` is
		// kept out of the picture the way SVG keeps it: as though it were
		// hidden, present in the drawing for something to point at and painted
		// nowhere. Building the nodes would paint it where it stands, and a
		// `<use>` would then paint it a second time.
		return n
	}
	for _, k := range e.Kids {
		n.Kids = append(n.Kids, img.node(k, st))
	}
	n.mask = maskUnder(def, st, n, warn)
	n.clip = clipUnder(n.clip, st, n, warn)
	return n
}
