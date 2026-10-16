package svg

import (
	"sync"

	"github.com/gabrielluizsf/antui/canvas"
)

// Image is a drawing read from SVG, ready to be painted at whatever size it is
// asked for. One is read once and painted as often as it is needed, which is
// what makes an icon cheaper than the file it came from after the first frame.
//
// An Image is safe to paint from more than one goroutine, and safe to paint
// while a picture a server is still fetching has yet to arrive: the cache of
// the last paint and the slot a picture lands in are each guarded, and a
// picture landing after a paint was made drops that canvas so the next paint
// shows it.
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

	// base is where the drawing was read from, which a relative address inside
	// it resolves against: the file [ParseFile] opened, and empty for a drawing
	// read from a string with no file of its own. An `xml:base` on an element
	// takes over from it for that element and everything under it. See
	// [resolveHref].
	base string

	warnings Warnings

	// warnMu guards warnings: a picture fetched from a server is said to have
	// failed by the goroutine that fetched it, while the drawing may be painted
	// and its warnings read at the same time.
	warnMu sync.Mutex

	// epoch counts the paints whose result may no longer be cached: it is
	// bumped when a picture arrives after a paint was made, so a paint already
	// running does not cache the hole it drew before it. See [Image.pictureArrived].
	epoch uint64

	// grads are the gradients the drawing declared, by the id a `url(#name)`
	// names. They are read before anything is built, because a shape may point
	// at a gradient written after it — which is the order a drawing is usually
	// written in, with the shapes first and the `<defs>` at the end.
	grads map[string]*gradient

	// pictures are the bitmaps the `<image>`s fetched, by the address each one
	// resolves to, so that a drawing pointing at the same file twice reads it
	// once — an animated GIF as every frame of its loop, chosen by the clock
	// when it is painted rather than kept here.
	// What is in here is never written to: an `<image>` that wants the picture
	// somewhere else or in another opacity is given where and how it fits rather
	// than changing what every other element pointing at the same address holds.
	pictures map[string]*pictureFrames

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
	// drawnTick is the frame every moving picture showed when the canvas above
	// was painted: the same tick means the same frames are still on the clock,
	// and another one means the canvas holds a frame that has ended.
	drawnTick int

	// animated says a picture that moves with the clock has landed in this
	// drawing, so which frame it shows is read from the clock as it is painted
	// and no canvas of it outlives the moment it was made in. Set when such a
	// picture arrives, under mu, beside the epoch it drops the cache with.
	animated bool
}
