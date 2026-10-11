package template

import "github.com/gabrielluizsf/antui/template/css"

// clampDim applies the min/max constraint of one axis to a container
// dimension, against the window like every other percentage.
func (c *CSS) clampDim(st css.Style, v int, axis string, base int) int {
	switch axis {
	case "width":
		if st.Has("max-width") && gridLengthDefinite(st.MaxWidth) {
			v = min(v, c.style.length(st, st.MaxWidth, base))
		}
		if st.Has("min-width") && gridLengthDefinite(st.MinWidth) {
			v = max(v, c.style.length(st, st.MinWidth, base))
		}
	case "height":
		if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) {
			v = min(v, c.style.length(st, st.MaxHeight, base))
		}
		if st.Has("min-height") && gridLengthDefinite(st.MinHeight) {
			v = max(v, c.style.length(st, st.MinHeight, base))
		}
	}
	return v
}
