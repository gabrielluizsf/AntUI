package canvas

import (
	"sync"

	"github.com/gabrielluizsf/antui/font"
)

// Face is a font at a size, ready to draw. Two kinds exist:
//
//	canvas.BuiltinFace()          the 8x16 bitmap, exactly as before
//	antui.SystemFace(15)          whatever this platform writes its own
//	                               interface in, at 15 pixels
//	canvas.LoadFace(path, 15)     a TrueType file of your own
//
// Set a face for the whole program:
//
//	canvas.SetDefaultFace(face)   // every Text, TextWidth and widget
type Face struct {
	// built-in, when font is nil; unit is how many device pixels one cell of
	// its 8x16 shape comes out at, so the font's own height is one and a size
	// between its steps is the shape with its cells at that size.
	unit  float64
	file  *font.TTF
	size  float64 // pixels per em
	scale float64 // font units to pixels

	ascent, descent, height int

	// next is the face handed a rune this one has no shape for — the rest
	// of a font-family list, then the face the program draws with by
	// default. It is set when the chain is made and never after, and is
	// read only by a face that has a file behind it: the bitmap font has no
	// outlines to fall back from, it is the end of every chain.
	next *Face

	// colors is the colour picture a rune draws here — nil for the ones
	// drawn as an outline — read off this file or off a face further down
	// the chain, and asked for once per rune.
	colors   map[rune]*font.ColorGlyph
	hasColor bool

	mu      sync.Mutex
	glyphs  map[rune]*font.Mask
	derived map[int]*Face
	// sized is the same face at sizes asked for outright rather than as a
	// multiple, keyed by the whole pixel it was made at. See [Face.AtSize].
	sized map[int]*Face
}
