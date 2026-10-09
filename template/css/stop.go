package css

import "github.com/gabrielluizsf/antui/canvas"

// Stop is one colour stop of a gradient. Offset is a fraction along the
// gradient line, or -1 when the declaration left it out and the painter is to
// spread it; Color may be [CurrentColor] until the cascade resolves it.
type Stop struct {
	Offset float64
	Color  canvas.Color
}
