package canvas

// growArea is the bounding box of a rectangle and a run of pixels on a row.
func growArea(a Area, x0, x1, y int) Area {
	if a.Width <= 0 {
		return Area{X: x0, Y: y, Width: x1 - x0 + 1, Height: 1}
	}
	if x0 < a.X {
		a.Width += a.X - x0
		a.X = x0
	}
	if x1 >= a.X+a.Width {
		a.Width = x1 - a.X + 1
	}
	if y < a.Y {
		a.Height += a.Y - y
		a.Y = y
	}
	if y >= a.Y+a.Height {
		a.Height = y - a.Y + 1
	}
	return a
}
