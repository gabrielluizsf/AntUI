package css

// gridAreaBounds folds each named cell's run over the rows into [min-column,
// min-row, max-column, max-row] and reports whether every area is rectangular
// and contiguous.
func gridAreaBounds(out [][]string) (bool, map[string][4]int) {
	bounds := make(map[string][4]int)
	for row, cells := range out {
		for col, name := range cells {
			if name == "." || name == "..." {
				continue
			}
			b, exists := bounds[name]
			if !exists {
				bounds[name] = [4]int{col, row, col, row}
				continue
			}
			b[0] = min(b[0], col)
			b[1] = min(b[1], row)
			b[2] = max(b[2], col)
			b[3] = max(b[3], row)
			bounds[name] = b
		}
	}
	for name, b := range bounds {
		for row := b[1]; row <= b[3]; row++ {
			for col := b[0]; col <= b[2]; col++ {
				if out[row][col] != name {
					return false, nil
				}
			}
		}
	}
	return true, bounds
}
