package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintPicture lays an `<image>`'s picture into the box it was given, through
// the transform that puts the drawing on the canvas — the same way writing is
// painted, with the element's own transform composed on at paint time rather
// than put into anything while the drawing is read.
//
// The picture's own pixels first go to where they belong inside the box: the
// fit says how big they come out and the alignment says where they sit when
// that leaves anything over, which is [picturePlacement]. A picture that fits
// goes straight onto the canvas; a sliced one is cut to the box on a picture of
// its own first, since what is cut is the rectangle the drawing asked for and
// nothing beyond it, however the picture was turned on the way.
func paintPicture(cv *canvas.Canvas, n *Node, m canvas.Matrix) {
	p := n.Pic
	// Which frame of a moving picture is showing is the clock's to say, read
	// here rather than kept in the slot: the frames are shared with every
	// element pointing at the same address, and the moment belongs to the paint.
	// The box comes with it for the reason [Picture.fitted] reads the two
	// together.
	pic, where := p.fitted(pictureClock())
	if pic == nil || pic.Width <= 0 || pic.Height <= 0 {
		return
	}
	bw, bh := where.box[2], where.box[3]
	if bw <= 0 || bh <= 0 {
		return
	}
	sx, sy, ox, oy := picturePlacement(where, float64(pic.Width), float64(pic.Height))
	total := m.Mul(n.Style.Transform)
	matrix := total.Mul(canvas.Translate(ox, oy).Mul(canvas.Scale(sx, sy)))
	if where.fit != PictureFitSlice {
		cv.BlitMatrix(pic, matrix, canvas.Area{Width: pic.Width, Height: pic.Height})
		return
	}
	layer, err := canvas.NewLayer(cv.Width, cv.Height)
	if err != nil {
		// There is no picture to cut: a canvas of no size has nothing to draw
		// on and nothing to draw into either, the same as every other layer
		// this package fails to make.
		return
	}
	layer.BlitMatrix(pic, matrix, sliceRegion(where.box, ox, oy, sx, sy, pic))
	layer.MaskShapes(canvas.MaskShape{Path: boxQuad(where.box, total), Rule: canvas.FillNonZero})
	cv.BlitOver(0, 0, layer)
}
