package svg

import "github.com/gabrielluizsf/antui/canvas"

// shapeBuild records the style inputs and the derived values made from them.
// Style.Dash is the effective pattern after pathLength, so it is kept here as
// a snapshot to tell an external edit from the package's own derived value.
type shapeBuild struct {
	ready         bool
	effectiveDash *canvas.Dash
	pathLength    float64
	hasStroke     bool
	width         float64
	transform     canvas.Matrix
	hasTransform  bool
	markerStart   *markerDef
	markerMid     *markerDef
	markerEnd     *markerDef
}

func sameShapeBuild(build shapeBuild, style Style) bool {
	return build.ready && sameDash(build.effectiveDash, style.Dash) && sameShapeInputs(build, style)
}

func sameShapeInputs(build shapeBuild, style Style) bool {
	return build.pathLength == style.pathLength && build.hasStroke == style.HasStroke &&
		build.width == style.Width && build.transform == style.Transform &&
		build.hasTransform == style.HasTransform && build.markerStart == style.markerStart &&
		build.markerMid == style.markerMid && build.markerEnd == style.markerEnd
}
