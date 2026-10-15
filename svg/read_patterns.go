package svg

import "github.com/gabrielluizsf/antui/canvas"

// readPatterns takes every `<pattern>` that named itself, under the id a `fill`
// or a `stroke` may point at, and builds the picture each one holds. They are
// named before any of them is built, the way the gradients and the masks are:
// a pattern may be written after the element that names it, and what is inside
// one may name anything else the drawing has — a gradient, a clip, a mask, a
// `<use>`, or another pattern written further down the file.
func (img *Image) readPatterns(root *element) {
	type declared struct {
		id string
		e  *element
	}
	var found []declared
	seen := map[string]bool{}
	var walk func(e *element)
	walk = func(e *element) {
		if e.Name == "pattern" {
			if id := e.attr("id"); id != "" && !seen[id] {
				seen[id] = true
				found = append(found, declared{id: id, e: e})
			}
		}
		for _, k := range e.Kids {
			walk(k)
		}
	}
	walk(root)

	img.patterns = make(map[string]*pattern, len(found))
	// Every pattern is named with an empty picture in it first, and only then
	// is any of them read and built: a pattern's picture may point at another
	// by an id that has not been reached yet, and a reference to a pattern that
	// is not there is a warning rather than a picture.
	for _, d := range found {
		img.patterns[d.id] = &pattern{transform: canvas.Identity()}
	}
	for _, d := range found {
		img.readPattern(d.id, d.e)
	}
	for _, d := range found {
		img.patterns[d.id].content = img.build(d.e, img.rootStyle(), true)
	}
}
