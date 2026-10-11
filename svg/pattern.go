package svg

import (
	"sync"

	"github.com/gabrielluizsf/antui/canvas"
)

// A `<pattern>` is a paint rather than a colour, the same as a gradient:
// `fill="url(#p)"` names one elsewhere in the drawing, and the shape is painted
// by asking the pattern what it has at each point rather than being painted one
// flat colour. What it has is a picture, written once, laid over the shape
// again and again from one edge of it to the other — which is what makes a
// chequered fill or a field of dots out of one small drawing.
//
// The picture is drawn into a canvas the size of one tile and kept there. Every
// pixel of the shape asks that canvas which part of the tile it is over, and
// the answer wraps round by the size of the tile, so the same picture comes
// back as many times as the shape is big enough to want it. The tile itself is
// drawn once for every shape that asks the same of it rather than once per
// shape, and once per pixel of neither, which is what keeps a pattern no dearer
// than a gradient to paint with.
//
// Three things say where the tile and its picture are, and each is read the way
// a gradient's geometry is: `patternUnits` says whether the tile is written in
// fractions of the box of the shape it paints or in the drawing's own
// coordinates — the box being the default — and `patternContentUnits` says the
// same for what is inside the tile, whose default is the other way round, the
// drawing's own. `patternTransform` moves the tile and everything in it
// without moving the shape it paints.
type pattern struct {
	// boxUnits is whether the tile below is written in fractions of the box of
	// the shape the pattern paints, which is what `patternUnits` chooses
	// between and what it says when it says nothing.
	boxUnits bool
	// contentBox is the same for the picture inside the tile, which is what
	// `patternContentUnits` chooses between — and there the drawing's own
	// coordinates are the default, not the fractions.
	contentBox bool
	// x, y, w, h is the tile, in whichever units boxUnits says: where it
	// starts and how big it is, before the pattern's own transform moves it.
	x, y, w, h float64
	// transform is `patternTransform`: it moves the tile and the picture in it,
	// the same as `gradientTransform` moves a gradient, and the shape the
	// pattern paints stays where it is. It is always a matrix that can be put
	// back, so an empty one is the identity rather than nothing.
	transform canvas.Matrix
	// content is the picture inside the tile, built where it was written — the
	// transform of the pattern and of anything inside it counts, because that
	// is part of how they were written — and drawn into the tile over and over
	// as the shape is painted. It is not part of the picture where it stands,
	// which is what forceKids says when it is built.
	content *Node
	// visible is `overflow`: whether the picture may run out of the tile. What
	// it does not say — which is what the spec's own stylesheet says a pattern
	// stands for — is that it may not, so the tile cuts what runs past its
	// edge unless `overflow="visible"` is written, and then the picture is
	// drawn into the tile from the copies of it that land there. See
	// [pattern.picture].
	visible bool
	// tiles is the picture this pattern already drew, kept for the shapes that
	// ask the same of it — the same tile at the same place and size, in the
	// same colour, over the same chain of patterns already being drawn, which
	// is the whole of what a picture is drawn from. Ten shapes under one
	// pattern cost one picture rather than ten, and drawn says how many were
	// drawn rather than found kept, which is what a test reads to know the
	// keeping holds — a picture reused and a picture drawn again come out the
	// same on the page. mu guards both because two goroutines may paint the
	// same drawing at once; the picture is drawn with mu off and only the
	// finding and the keeping of it are under mu, so a pattern inside a pattern
	// never waits on the one holding it — and where two draw the same picture
	// at once, whichever keeps it first stands and the other's is let go rather
	// than both being held.
	mu    sync.Mutex
	tiles []tilePicture
	drawn int
}
