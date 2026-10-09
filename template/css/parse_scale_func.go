package css

// parseScaleFunc reads scale(), scaleX() and scaleY().
func parseScaleFunc(name string, args []string) (TransformFunc, bool) {
	switch name {
	case "scale":
		if len(args) != 1 && len(args) != 2 {
			return TransformFunc{}, false
		}
		sx, ok := transformFactor(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		sy := sx
		if len(args) == 2 {
			if sy, ok = transformFactor(args[1]); !ok {
				return TransformFunc{}, false
			}
		}
		return TransformFunc{Kind: TransformScale, Sx: sx, Sy: sy}, true
	case "scalex":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		sx, ok := transformFactor(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformScale, Sx: sx, Sy: 1}, true
	case "scaley":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		sy, ok := transformFactor(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformScale, Sx: 1, Sy: sy}, true
	}
	return TransformFunc{}, false
}
