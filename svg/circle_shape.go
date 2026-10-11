package svg

import "github.com/gabrielluizsf/antui/canvas"

// circleShape is a circle, which is an ellipse with one radius.
func circleShape(e *element, warn func(string, ...any)) *canvas.Path {
	cx, _ := parseLength(e.attr("cx"))
	cy, _ := parseLength(e.attr("cy"))
	r, ok := parseLength(e.attr("r"))
	if !ok || r <= 0 {
		return nil
	}
	p := canvas.NewPath()
	p.AddCircle(cx, cy, r)
	return p
}
