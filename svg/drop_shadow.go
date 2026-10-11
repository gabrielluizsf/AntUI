package svg

import "github.com/gabrielluizsf/antui/canvas"

// dropShadow is a `drop-shadow(...)`: the offset, the blur radius and the
// colour the silhouette of the element's picture takes. A colour that was never
// written is `currentColor`, which is why the flag is kept rather than a colour
// resolved here — only [applyFilters] is given the colour of the widget.
type dropShadow struct {
	dx, dy, blur float64
	colour       canvas.Color
	current      bool
}
