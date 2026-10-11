package svg

import "github.com/gabrielluizsf/antui/canvas"

// readClipPaths takes every `<clipPath>` that named itself, under the id a
// `clip-path` may point at, and works out the shapes it holds. One written
// without an id cannot be pointed at, so it is left unread the way a gradient
// with no id is: nothing in the drawing can reach it.
func (img *Image) readClipPaths(root *element) {
	type declared struct {
		id   string
		e    *element
		rule canvas.FillRule
	}
	var found []declared
	seen := map[string]bool{}
	// The clip-rule in force walks down with the elements, because it is
	// inherited and a clipPath may stand anywhere in the drawing — inside a
	// group that wrote it once for everything below it. It is read on the way
	// down rather than by building the elements above the clipPath, because
	// those are part of the picture too and every one of them is built and
	// read once already: reading them again here would say everything twice.
	var walk func(e *element, rule canvas.FillRule)
	walk = func(e *element, rule canvas.FillRule) {
		rule = clipRuleIn(e, rule)
		if e.Name == "clipPath" {
			if id := e.attr("id"); id != "" && !seen[id] {
				seen[id] = true
				found = append(found, declared{id: id, e: e, rule: rule})
			}
		}
		for _, k := range e.Kids {
			walk(k, rule)
		}
	}
	walk(root, canvas.FillNonZero)

	img.clips = make(map[string]*clipPath, len(found))
	// Every clipPath is named before any of them is built, so that one
	// written inside another — or pointed at from inside one — finds a place
	// to land rather than finding nothing there at all.
	for _, d := range found {
		img.clips[d.id] = &clipPath{}
	}
	for _, d := range found {
		img.clips[d.id] = img.readClipPath(d.e, d.rule)
	}
	// Every clipPath holds its picture now, and every `clip-path` written
	// inside one of them is still a name pointing at a clip whose shapes are
	// not worked out yet. They are worked out here, one id at a time,
	// following what is inside each first: a clip written inside a clipPath
	// cuts the shape it is written on, and it may name one written further
	// down the file, which can only be followed now that every id is known
	// and every picture is built. The walk goes round a circle at the first
	// repeat rather than for ever — see [Image.measureClip].
	gray := map[string]bool{}
	for _, d := range found {
		img.measureClip(d.id, gray)
	}
}
