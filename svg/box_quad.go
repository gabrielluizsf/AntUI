package svg

import "github.com/gabrielluizsf/antui/canvas"

// boxQuad is the box of a picture as the shape its slice is cut to: the four
// corners taken through the transform that puts the drawing on the canvas and
// joined into one outline. The axis-aligned rectangle around those corners
// (what [regionArea] answers) is bigger than the box whenever the drawing was
// turned, and cutting to it leaves the picture sticking out past its own box
// on every rotated edge — the blit reaches about a source pixel further than
// the box on each side, which is exactly what the cut is there to remove. The
// box itself is the truth however it was turned.
func boxQuad(box [4]float64, m canvas.Matrix) *canvas.Path {
	x, y, w, h := box[0], box[1], box[2], box[3]
	ax, ay := m.Map(x, y)
	bx, by := m.Map(x+w, y)
	cx, cy := m.Map(x+w, y+h)
	dx, dy := m.Map(x, y+h)
	p := canvas.NewPath()
	p.MoveTo(ax, ay)
	p.LineTo(bx, by)
	p.LineTo(cx, cy)
	p.LineTo(dx, dy)
	p.Close()
	return p
}
