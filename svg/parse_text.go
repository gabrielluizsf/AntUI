package svg

// parseText reads a drawing that came from path — "" when it came from no file
// at all — and keeps path as the base a relative address inside it resolves
// against, which is what every `<image>` later asks for while it is built.
// chain is the drawings read through one another to get here, which a picture
// pointing back at one of them is stopped by.
func parseText(src, path string, chain []string) (*Image, error) {
	root, err := parseDocument(src)
	if err != nil {
		return nil, err
	}
	if root.foreign {
		return nil, &Error{What: "the root tag stands in a namespace other than the SVG one"}
	}
	if root.Name != "svg" {
		return nil, &Error{What: "the root tag is <" + root.Name + ">, not <svg>"}
	}
	img := imageFor(root, path, chain)
	img.readSize(root)
	img.readDocFit(root)
	// The gradients come before the shapes, because a shape may name one that is
	// only written later in the file, and every gradient that names another has
	// to be followed before either is asked about a colour.
	img.readGradients(root)
	// The ids come with them, for the same reason: a `<use>` may name an element
	// that is written after it, and only the whole drawing can say what an id is.
	img.readIds(root)
	// The clip paths come last of the things that are pointed at, since a
	// `clip-path` may name one written after the element it cuts, and building
	// them is building shapes that may themselves be pointed at.
	img.readClipPaths(root)
	// The filters come with them: what is inside a `<filter>` is only what its
	// primitives say — no picture painted with the rest of the drawing — so the
	// ids being in place is all one of them needs to be read, and an element's
	// `filter="url(#name)"` may name one written further down.
	img.readFilterDefs(root)
	// The masks come after them, for the same reason taken once more: a mask's
	// picture may be cut by a clip, filled with a gradient, drawn through a
	// `<use>` and named from inside another mask, so every id the drawing has
	// to offer has to be known before one of them is built.
	img.readMasks(root)
	// The patterns come last of the things that are pointed at: what is inside
	// one is a picture painted with everything else — a gradient, a clip, a
	// mask, a `<use>`, another pattern — so all of it has to be in place first.
	img.readPatterns(root)
	// The circles the masks and the patterns close are walked now, while every
	// picture they hold is built and before any of it is painted: what closes
	// a circle is said here, once, rather than only left off in silence while
	// the drawing is being painted.
	img.warnMaskPatternCycles()
	// The markers come after them, and are last for the same reason taken once
	// more: what is inside a marker is a picture painted with everything else,
	// and a marker may name another written further down the file.
	img.readMarkers(root)
	img.read(root)
	return img, nil
}
