package css

// parseSkewFunc reads skew(), skewX() and skewY().
func parseSkewFunc(name string, args []string) (TransformFunc, bool) {
	switch name {
	case "skew":
		if len(args) != 1 && len(args) != 2 {
			return TransformFunc{}, false
		}
		ax, ok := ParseAngle(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		ay, _ := ParseAngle(args[1])
		return TransformFunc{Kind: TransformSkew, Ax: ax, Ay: ay}, true
	case "skewx":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		ax, ok := ParseAngle(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformSkew, Ax: ax}, true
	case "skewy":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		ay, ok := ParseAngle(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformSkew, Ay: ay}, true
	}
	return TransformFunc{}, false
}

// parseMatrixFunc reads the six numbers of matrix().
func parseMatrixFunc(args []string) (TransformFunc, bool) {
	if len(args) != 6 {
		return TransformFunc{}, false
	}
	var m [6]float64
	for i, a := range args {
		v, ok := transformNumber(a)
		if !ok {
			return TransformFunc{}, false
		}
		m[i] = v
	}
	return TransformFunc{Kind: TransformMatrix, M: m}, true
}
