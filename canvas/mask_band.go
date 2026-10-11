package canvas

import "math"

// maskBand is the box every shape can reach, clamped to the clip region, and
// whether it holds anything at all.
func maskBand(shapes []MaskShape, cx0, cy0, cx1, cy1 int) (int, int, int, int, bool) {
	bx0, by0, bx1, by1 := cx1, cy1, cx0, cy0
	for _, s := range shapes {
		if s.Path == nil {
			continue
		}
		minX, minY, maxX, maxY, ok := s.Path.Bounds()
		if !ok {
			continue
		}
		x0 := clampInt(int(math.Floor(minX)), cx0, cx1)
		y0 := clampInt(int(math.Floor(minY)), cy0, cy1)
		x1 := clampInt(int(math.Ceil(maxX)), cx0, cx1)
		y1 := clampInt(int(math.Ceil(maxY)), cy0, cy1)
		if x0 >= x1 || y0 >= y1 {
			continue
		}
		bx0, by0 = min(bx0, x0), min(by0, y0)
		bx1, by1 = max(bx1, x1), max(by1, y1)
	}
	return bx0, by0, bx1, by1, bx0 < bx1 && by0 < by1
}
