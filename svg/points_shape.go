package svg

import "github.com/gabrielluizsf/antui/canvas"

// pointsShape is a polyline or a polygon: a list of points, which a polygon
// joins back to where it started and a polyline leaves open.
func pointsShape(s string, close bool, warn func(string, ...any)) *canvas.Path {
	fields := splitArgs(s)
	// The points may be written as a flat list of numbers, two to a point, or
	// as a list of pairs separated by anything else.
	nums := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, err := parseNumber(f)
		if err != nil {
			warn("points needs a list of numbers, and %q is not one", f)
			return nil
		}
		nums = append(nums, v)
	}
	if len(nums)%2 != 0 {
		warn("points needs two numbers to a point, and there are %d of them", len(nums))
		return nil
	}
	if len(nums) < 4 {
		// A single point is not a line: there is no line to draw and no area to
		// fill, so there is nothing here to paint.
		return nil
	}
	pts := make([]canvas.Point, 0, len(nums)/2)
	for i := 0; i < len(nums); i += 2 {
		pts = append(pts, canvas.Point{X: nums[i], Y: nums[i+1]})
	}
	p := canvas.NewPath()
	p.AddPolyline(pts, close)
	return p
}
