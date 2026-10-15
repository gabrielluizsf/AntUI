package svg

import "github.com/gabrielluizsf/antui/canvas"

// ellipseShape is an ellipse about its own centre.
func ellipseShape(e *element, warn func(string, ...any)) *canvas.Path {
	cx, _ := parseLength(e.attr("cx"))
	cy, _ := parseLength(e.attr("cy"))
	rx, okX := parseLength(e.attr("rx"))
	ry, okY := parseLength(e.attr("ry"))
	if !okX || !okY || rx <= 0 || ry <= 0 {
		return nil
	}
	p := canvas.NewPath()
	p.AddEllipse(canvas.Point{X: cx, Y: cy}, canvas.Point{X: rx, Y: ry})
	return p
}
