package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// boxSide is one side of one box, measured: which side it is (a
// [css.SideTop]/[css.SideRight]/[css.SideBottom]/[css.SideLeft]), the box it
// belongs to and where it sits on the screen, how thick the stroke is, and the
// colour to lay it in.
//
// Every border style is painted the same way around — solid, dashed, dotted,
// double all take this and nothing else — so the box is carried whole rather
// than as eight numbers each painter would otherwise take and unpack again.
type boxSide struct {
	side   int
	x, y   int
	w, h   int
	width  int // the stroke's thickness
	colour canvas.Color
}

// horizontal says whether the side runs left to right rather than top to
// bottom, which is what decides how far along it there is to paint.
func (s boxSide) horizontal() bool {
	return s.side == css.SideTop || s.side == css.SideBottom
}

// along is how far the side runs: the box's width when it is horizontal, its
// height when it is not.
func (s boxSide) along() int {
	if s.horizontal() {
		return s.w
	}
	return s.h
}

// rect is a rectangle covering part of the side: from `at` to `at+n` along it,
// `thick` across it, sitting `off` in from the box's edge that this side is
// flush with. Every border style paints nothing else — a solid side is one of
// these spanning the whole run, a dash is a short one further along, and one
// line of a double is one further in from the edge.
func (s boxSide) rect(at, n, off, thick int) (x, y, w, h int) {
	if s.horizontal() {
		y := s.y + off
		if s.side == css.SideBottom {
			y = s.y + s.h - off - thick
		}
		return s.x + at, y, n, thick
	}
	x = s.x + off
	if s.side == css.SideRight {
		x = s.x + s.w - off - thick
	}
	return x, s.y + at, thick, n
}
