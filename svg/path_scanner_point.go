package svg

import "github.com/gabrielluizsf/antui/canvas"

// point reads the two numbers of a point.
func (s *pathScanner) point() (canvas.Point, bool) {
	x, ok1 := s.number()
	y, ok2 := s.number()
	return canvas.Point{X: x, Y: y}, ok1 && ok2
}
