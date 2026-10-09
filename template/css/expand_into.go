package css

func expandInto(out [4]Length, vals []Length) [4]Length {
	switch len(vals) {
	case 1:
		return [4]Length{vals[0], vals[0], vals[0], vals[0]}
	case 2:
		return [4]Length{vals[0], vals[1], vals[0], vals[1]}
	case 3:
		return [4]Length{vals[0], vals[1], vals[2], vals[1]}
	default:
		return [4]Length{vals[0], vals[1], vals[2], vals[3]}
	}
}
