package font

// transformContour places one component's outline under its transform.
func transformContour(comp component, part Contour) Contour {
	moved := Contour{
		X:  make([]float64, len(part.X)),
		Y:  make([]float64, len(part.Y)),
		On: part.On,
	}
	for i := range part.X {
		moved.X[i] = comp.a*part.X[i] + comp.c*part.Y[i] + comp.dx
		moved.Y[i] = comp.b*part.X[i] + comp.d*part.Y[i] + comp.dy
	}
	return moved
}
