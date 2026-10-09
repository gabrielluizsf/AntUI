package css

// parseTranslateFunc reads translate(), translateX() and translateY().
func parseTranslateFunc(name string, args []string, ctx Units) (TransformFunc, bool) {
	switch name {
	case "translate":
		if len(args) != 1 && len(args) != 2 {
			return TransformFunc{}, false
		}
		dx, ok := transformLength(args[0], ctx)
		if !ok {
			return TransformFunc{}, false
		}
		dy := Zero()
		if len(args) == 2 {
			if dy, ok = transformLength(args[1], ctx); !ok {
				return TransformFunc{}, false
			}
		}
		return TransformFunc{Kind: TransformTranslate, Dx: dx, Dy: dy}, true
	case "translatex":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		dx, ok := transformLength(args[0], ctx)
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformTranslate, Dx: dx}, true
	case "translatey":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		dy, ok := transformLength(args[0], ctx)
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformTranslate, Dy: dy}, true
	}
	return TransformFunc{}, false
}
