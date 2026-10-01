// Package svg reads a drawing written in the Scalable Vector Graphics language
// and paints it on a canvas, so an icon can live in its own file and be used the
// way a font glyph is: named once, drawn wherever it is asked for, at whatever
// size it is given.
//
// The drawing is read into a tree of nodes that keeps the attributes it was
// written with, and is painted by walking that tree with the style a node
// inherits from the one above it. What is left out is said out loud rather than
// dropped in silence: an image that quietly loses half of itself is harder to
// track down than one that says which half it could not draw.
package svg

import (
	"fmt"
	"strings"
)

// Warning is one thing a drawing asked for that this package could not do. They
// are collected while the drawing is walked and can be read back after it is
// painted, which is how a caller finds out that an icon came out short of what
// its author drew without having to read the file again.
type Warning struct {
	// Element is the tag the warning belongs to, without its namespace.
	Element string
	// What says what was wanted, in a form a person reading a warning can act
	// on, such as `stroke-dasharray needs a list of lengths`.
	What string
}

// String is the warning on one line.
func (w Warning) String() string { return w.Element + ": " + w.What }

// Warnings is a list of them, kept in the order they were found.
type Warnings []Warning

// String is the whole list, one warning per line, or the words for it having
// none.
func (ws Warnings) String() string {
	if len(ws) == 0 {
		return "no warnings"
	}
	out := make([]string, len(ws))
	for i, w := range ws {
		out[i] = w.String()
	}
	return strings.Join(out, "\n")
}

// warn records one thing that could not be drawn, naming the element it came
// from so a reader can go and look at it.
func (ws *Warnings) warn(element, format string, args ...any) {
	*ws = append(*ws, Warning{Element: element, What: fmt.Sprintf(format, args...)})
}
