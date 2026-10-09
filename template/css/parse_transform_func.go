package css

// parseTransformFunc reads one function's arguments. translateX/Y and skewX/Y
// fold into their two-axis siblings, with the missing axis at its zero.
func parseTransformFunc(name, inner string, ctx Units) (TransformFunc, bool) {
	args, ok := transformArgs(inner)
	if !ok {
		return TransformFunc{}, false
	}
	switch name {
	case "translate", "translatex", "translatey":
		return parseTranslateFunc(name, args, ctx)
	case "scale", "scalex", "scaley":
		return parseScaleFunc(name, args)
	case "rotate":
		return parseRotateFunc(args)
	case "skew", "skewx", "skewy":
		return parseSkewFunc(name, args)
	case "matrix":
		return parseMatrixFunc(args)
	}
	return TransformFunc{}, false
}
