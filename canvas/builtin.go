package canvas

// builtin is the 8x16 bitmap, as a Face.
var builtin = &Face{
	unit:    1,
	ascent:  12,
	descent: FontHeight - 12,
	height:  FontHeight,
}
