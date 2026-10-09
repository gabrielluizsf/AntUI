package css

// parseRotateFunc reads rotate().
func parseRotateFunc(args []string) (TransformFunc, bool) {
	if len(args) != 1 {
		return TransformFunc{}, false
	}
	a, ok := ParseAngle(args[0])
	if !ok {
		return TransformFunc{}, false
	}
	return TransformFunc{Kind: TransformRotate, Ax: a}, true
}
