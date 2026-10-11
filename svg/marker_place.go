package svg

import "github.com/gabrielluizsf/antui/canvas"

// markerPlace is where one marker sits on one shape: the marker itself and the
// matrix that takes its picture to the drawing's coordinates, with the element's
// own transform already in it under the same condition the shape's path got it,
// so that a marker follows the vertex the way the vertex follows the shape.
type markerPlace struct {
	def *markerDef
	at  canvas.Matrix
	// room is the rectangle of the room the marker asked for, taken back into
	// the coordinates its picture was written in so that it cuts the picture
	// where the room's edges land on it. It is nil where `overflow` says the
	// picture may run out of the room, which draws it whole — and nil where
	// the fit of a viewBox has no way back, which collapses the picture
	// anyway rather than cutting it the wrong way.
	room *clipRegion
	// prop is the property that asked for this one — "marker-start",
	// "marker-mid" or "marker-end" — and is only here so that a warning about
	// a marker going round for ever can say which of the three it was.
	prop string
}
