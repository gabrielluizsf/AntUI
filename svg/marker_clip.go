package svg

import "github.com/gabrielluizsf/antui/canvas"

// markerClip is the room as a cut, in the coordinates the picture is written
// in: its rectangle taken backwards through the fit of the viewBox, so that
// where the room's edges land on the picture follow where the picture's own
// writing lands. A visible overflow gets no cut, and a fit with no way back
// collapses the picture to nothing anyway, so it gets no cut rather than the
// wrong one.
func markerClip(onto canvas.Matrix, w, h float64, visible bool) *clipRegion {
	if visible {
		return nil
	}
	inv, ok := onto.Inverse()
	if !ok {
		return nil
	}
	p := canvas.NewPath()
	p.AddRect(canvas.Point{X:0,Y:0}, canvas.Point{X:w,Y:h})
	p.Transform(inv)
	return &clipRegion{shapes: []canvas.MaskShape{{Path: p, Rule: canvas.FillNonZero}}}
}
