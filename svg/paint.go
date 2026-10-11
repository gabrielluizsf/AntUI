package svg

import "github.com/gabrielluizsf/antui/canvas"

// paint is the drawing itself: a canvas twice the size asked for, the nodes
// walked in order with the transform that maps the drawing's own coordinates
// onto it, and the result brought back down to size.
func (img *Image) paint(width, height int, current canvas.Color) *canvas.Canvas {
	// A layer rather than a plain canvas, because an icon is a picture with soft
	// edges that gets laid over a widget, not a surface of its own. A plain
	// canvas composites over what is there and answers opaque, which on an empty
	// one means a quarter-covered edge is stored as a quarter-bright opaque
	// pixel: a dark rim all the way round the icon.
	big, err := canvas.NewLayer(width*supersample, height*supersample)
	if err != nil {
		return nil
	}
	// The viewBox maps onto the whole of the drawing, and the drawing maps onto
	// the canvas, so the two compose into the transform every node is drawn
	// with. A drawing with a viewBox that does not start at the origin is
	// shifted, which is what keeps a viewBox of "0 0 24 24" and one of
	// "-12 -12 24 24" both filling the icon the same way.
	m := fitTransform(img.ViewBox, float64(width*supersample), float64(height*supersample))
	paintNodes(big, img.Root, m, current, nil, nil)

	if supersample == 1 {
		return big
	}
	// Scaled averages the block each pixel covers weighted by alpha and writes
	// the result rather than compositing it, so what comes down keeps the
	// transparency the drawing had. Halving by Blit would blend instead and
	// flatten every edge to solid.
	return canvas.Scaled(big, width, height)
}
