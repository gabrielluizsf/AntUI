package svg

import "github.com/gabrielluizsf/antui/canvas"

// lineShape is a straight line from one point to the other. It is drawn only by
// its stroke: SVG fills nothing, because a line encloses no area, so a line
// with a fill on it shows its fill nowhere.
func lineShape(e *element) *canvas.Path {
	x1, ok1 := parseLength(e.attr("x1"))
	y1, ok1b := parseLength(e.attr("y1"))
	x2, ok2 := parseLength(e.attr("x2"))
	y2, ok2b := parseLength(e.attr("y2"))
	if !ok1 || !ok1b || !ok2 || !ok2b {
		return nil
	}
	p := canvas.NewPath()
	p.MoveTo(x1, y1)
	p.LineTo(x2, y2)
	return p
}
