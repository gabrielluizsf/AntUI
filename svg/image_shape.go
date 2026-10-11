package svg

import "github.com/gabrielluizsf/antui/canvas"

// shape is the outline a node draws, or nil when it draws none: a group, a
// `<defs>` of things nothing points at, or a tag this package does not know.
func (img *Image) shape(e *element, st Style) *canvas.Path {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	switch e.Name {
	case "path":
		return pathData(e.attr("d"), warn)
	case "rect":
		return rectShape(e, warn)
	case "circle":
		return circleShape(e, warn)
	case "ellipse":
		return ellipseShape(e, warn)
	case "line":
		return lineShape(e)
	case "polyline":
		return pointsShape(e.attr("points"), false, warn)
	case "polygon":
		return pointsShape(e.attr("points"), true, warn)
	}
	switch e.Name {
	case "g", "svg", "a", "switch":
		return nil
	}
	if isKnownButUnpainted(e.Name) {
		return nil
	}
	warn("%q is not a shape this package draws, so it is left out", e.Name)
	return nil
}
