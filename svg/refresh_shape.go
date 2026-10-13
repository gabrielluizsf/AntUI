package svg

import "github.com/gabrielluizsf/antui/canvas"

// refreshShape rebuilds values derived from a shape's mutable style. The
// original path is used every time so a changed dash is cut before the
// transform, and a changed marker is placed on the original vertex.
func refreshShape(n *Node) bool {
	if n == nil || n.sourcePath == nil {
		return false
	}
	if n.shapeBuild.ready && sameDash(n.Style.Dash, n.shapeBuild.effectiveDash) &&
		sameShapeInputs(n.shapeBuild, n.Style) {
		return false
	}
	rawDash := n.rawDash
	if !n.shapeBuild.ready || !sameDash(n.Style.Dash, n.shapeBuild.effectiveDash) {
		rawDash = cloneDash(n.Style.Dash)
	}
	style := n.Style
	style.Dash = cloneDash(rawDash)
	applyPathLength(&style, n.sourcePath)
	if sameShapeBuild(n.shapeBuild, style) {
		n.rawDash = rawDash
		return false
	}

	n.rawDash = rawDash
	n.Style = style
	n.Path = clonePath(n.sourcePath)
	if style.HasTransform && style.Transform != (canvas.Matrix{}) {
		n.Path.Transform(style.Transform)
	}
	n.dashed = nil
	if style.HasStroke && style.Width > 0 && style.Dash != nil {
		if cut := n.sourcePath.Dashed(*style.Dash); cut != n.sourcePath {
			n.dashed = clonePath(cut)
			if style.HasTransform && style.Transform != (canvas.Matrix{}) {
				n.dashed.Transform(style.Transform)
			}
		}
	}
	n.markers = placeMarkers(n.sourcePath, style)
	n.shapeBuild = shapeBuild{
		ready: true, effectiveDash: cloneDash(style.Dash), pathLength: style.pathLength,
		hasStroke: style.HasStroke, width: style.Width, transform: style.Transform,
		hasTransform: style.HasTransform, markerStart: style.markerStart,
		markerMid: style.markerMid, markerEnd: style.markerEnd,
	}
	return true
}
