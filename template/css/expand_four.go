package css

import "github.com/gabrielluizsf/antui/canvas"

// The three expandFour variants repeat a one-to-four CSS value across the
// four sides: one value fills all, two alternate, three leave out the fourth
// (it mirrors the second), four pass through.
func expandFour(v []int) [4]int {
	switch len(v) {
	case 1:
		return [4]int{v[0], v[0], v[0], v[0]}
	case 2:
		return [4]int{v[0], v[1], v[0], v[1]}
	case 3:
		return [4]int{v[0], v[1], v[2], v[1]}
	}
	return [4]int{v[0], v[1], v[2], v[3]}
}

func expandFourInt8(v []uint8) [4]uint8 {
	switch len(v) {
	case 1:
		return [4]uint8{v[0], v[0], v[0], v[0]}
	case 2:
		return [4]uint8{v[0], v[1], v[0], v[1]}
	case 3:
		return [4]uint8{v[0], v[1], v[2], v[1]}
	}
	return [4]uint8{v[0], v[1], v[2], v[3]}
}

func expandFourColor(v []canvas.Color) [4]canvas.Color {
	switch len(v) {
	case 1:
		return [4]canvas.Color{v[0], v[0], v[0], v[0]}
	case 2:
		return [4]canvas.Color{v[0], v[1], v[0], v[1]}
	case 3:
		return [4]canvas.Color{v[0], v[1], v[2], v[1]}
	}
	return [4]canvas.Color{v[0], v[1], v[2], v[3]}
}
