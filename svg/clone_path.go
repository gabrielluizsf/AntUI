package svg

import "github.com/gabrielluizsf/antui/canvas"

// clonePath is a copy of a path, because the one on the node is shared by every
// frame that paints this drawing and must not be transformed in place.
func clonePath(p *canvas.Path) *canvas.Path {
	if p == nil {
		return nil
	}
	pts, closed := p.Points()
	out := canvas.NewPath()
	for i, sub := range pts {
		if len(sub) == 0 {
			continue
		}
		out.MoveTo(sub[0].X, sub[0].Y)
		for _, q := range sub[1:] {
			out.LineTo(q.X, q.Y)
		}
		if i < len(closed) && closed[i] {
			out.Close()
		}
	}
	return out
}
