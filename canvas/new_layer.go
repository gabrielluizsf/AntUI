package canvas

// NewLayer makes an ARGB32 canvas that keeps what is drawn into it transparent
// where nothing has been drawn. It is for a picture that is assembled out of
// shapes with soft edges and then laid over something else — an icon, a
// shadow — rather than for a surface of its own.
//
// A plain canvas composites every colour over what is there and answers opaque,
// because that is right for a window's front buffer and fast: the destination
// is known to be solid, so the result of drawing on it always is. A layer has
// no such destination, and the same arithmetic on it multiplies the colour by
// the coverage against nothing and stores that as if it were the colour: a
// shape whose edge is a quarter covered comes out a quarter as bright and
// fully opaque, which is a dark rim around every soft edge rather than a soft
// edge. So a layer composites with [BlendOver], which keeps both alphas.
//
// Everything that draws works on a layer — the shapes, the paths, the strokes —
// and the drawing is the same arithmetic one pixel at a time, so a layer is for
// pictures built once and drawn many times rather than for a per-frame surface.
func NewLayer(width, height int) (*Canvas, error) {
	cv, err := NewCanvas(width, height)
	if err != nil {
		return nil, err
	}
	cv.over = true
	return cv, nil
}
