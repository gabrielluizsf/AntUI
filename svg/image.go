package svg

import (
	"sync"

	"github.com/gabrielluizsf/antui/canvas"
)

// Image is a drawing read from SVG, ready to be painted at whatever size it is
// asked for. One is read once and painted as often as it is needed, which is
// what makes an icon cheaper than the file it came from after the first frame.
//
// An Image is safe to paint from more than one goroutine as long as nothing
// writes to it: painting caches the last size it was drawn at, and that cache
// is guarded.
type Image struct {
	// Width and Height are the size the drawing says it is, in the units of its
	// own viewBox, which is what a caller asking for an icon at twenty-four
	// pixels wants to know it was drawn to be.
	Width, Height float64
	// ViewBox is the rectangle of the drawing's own coordinates that maps onto
	// the size above, as minX, minY, width, height. A drawing with no viewBox
	// has one covering the whole of it.
	ViewBox [4]float64
	// Root is the top of the drawing: the `<svg>` and everything inside it.
	Root *Node

	warnings Warnings

	// grads are the gradients the drawing declared, by the id a `url(#name)`
	// names. They are read before anything is built, because a shape may point
	// at a gradient written after it — which is the order a drawing is usually
	// written in, with the shapes first and the `<defs>` at the end.
	grads map[string]*gradient

	// pictures are the bitmaps the `<image>`s fetched, by the address each one
	// named, so that a drawing pointing at the same file twice reads it once.
	// What is in here is never written to: an `<image>` that wants the picture
	// somewhere else or in another opacity is given where and how it fits rather
	// than changing what every other element pointing at the same address holds.
	pictures map[string]*canvas.Canvas

	// clips are the clip paths the drawing declared, by the id a `clip-path`
	// names. They are read before anything is built for the same reason as the
	// gradients: a clip is written at the end of the drawing beside them, and
	// what points at it is written first.
	clips map[string]*clipPath

	// masks are the masks the drawing declared, by the id a `mask` names, read
	// for the same reason as the clips and the gradients: a mask is written at
	// the end of the drawing beside them, and what points at it is written
	// first — sometimes inside another mask.
	masks map[string]*maskDef

	// filters are the `<filter>` elements the drawing declared, by the id an
	// element's `filter="url(#name)"` names, read after the ids are all in
	// place and before anything is built, for the same reason as the rest of
	// what is pointed at: one may be written after the element that points at
	// it. What is inside a filter is only what its primitives say — no picture
	// painted with the rest of the drawing — so nothing else has to be in
	// place for one to be read. See [Image.readFilterDefs].
	filters map[string]*filterDef

	// patterns are the patterns the drawing declared, by the id a `fill` or a
	// `stroke` names, read after the masks and for the same reason once more:
	// what is inside a pattern paints with everything the drawing has, and a
	// reference to one written further down the file would find nothing there.
	patterns map[string]*pattern

	// markers are the markers the drawing declared, by the id a `marker-start`
	// names, read after the patterns and for the same reason yet again: what is
	// inside a marker is a picture painted with everything else — a gradient, a
	// clip, a mask, a `<use>`, another marker — so all of it has to be in place
	// before one of them is built.
	markers map[string]*markerDef

	// ids are the elements that wrote an id, under that id, so that a `<use>`
	// can be followed to the element it names — which may be written after the
	// `<use>`, and usually is, at the end of the drawing next to the gradients.
	ids map[string]*element

	// uses are the ids being drawn through a `<use>` right now, outermost
	// first, while the drawing is being read. A reference that names one of
	// them closes a circle rather than drawing a shape, so it is stopped at
	// the first repeat.
	uses []string

	// The last size and colour painted and the canvas it was painted on, so that
	// a frame which asks for the same icon again costs nothing. Guarded by mu,
	// because two goroutines may ask for the same icon at once and the worst
	// that can happen is that one of them draws it.
	mu          sync.Mutex
	drawn       *canvas.Canvas
	drawnWidth  int
	drawnHeight int
	// drawnColour is the `currentColor` the cached canvas was painted in, which
	// is part of what makes it the canvas asked for: an icon painted red is not
	// the icon asked for a moment later in blue.
	drawnColour canvas.Color
}
