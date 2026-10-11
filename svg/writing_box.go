package svg

import "strings"

// writingBox is the box the writing covers in the drawing's own coordinates,
// measured where the letters are before the transform moves them.
func writingBox(pieces []piece) [4]float64 {
	var b [4]float64
	found := false
	add := func(x0, y0, x1, y1 float64) {
		if !found {
			b = [4]float64{x0, y0, x1, y1}
			found = true
			return
		}
		b[0], b[1] = min(b[0], x0), min(b[1], y0)
		b[2], b[3] = max(b[2], x1), max(b[3], y1)
	}
	for _, p := range pieces {
		if p.run.Path != nil {
			// The piece is hung along an outline rather than at a pen, so what
			// it covers is where the outline is rather than where the pen was.
			if x0, y0, x1, y1, ok := p.run.Path.Bounds(); ok {
				add(x0, y0, x1, y1)
			}
			continue
		}
		ascent := float64(p.face.Ascent()) / p.scale
		descent := float64(p.face.Descent()) / p.scale
		// A piece with a line break in it is as many lines tall as it has
		// breaks, and the ones below the first are below the pen.
		down := float64(1+strings.Count(p.run.Text, "\n")) *
			float64(p.face.Height()) / p.scale
		add(p.penX, p.penY-ascent, p.penX+p.width, p.penY+descent+down)
	}
	if !found {
		return [4]float64{}
	}
	return b
}
