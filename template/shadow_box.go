package template

import "github.com/gabrielluizsf/antui/template/css"

// shadowBox is the box a shadow is cast from, and the face it falls on: the
// rectangle every shadow of one element shares, the corner radii it is cut to,
// and whether the shadow falls inside the box rather than behind it.
//
// Shadows are painted in two passes around the box — the outer ones before the
// surface goes down and the inner ones after — so the inset flag is part of what
// is being asked for rather than something each call spells out.
type shadowBox struct {
	st         css.Style
	x, y, w, h int
	rx, ry     int
	inset      bool
}
