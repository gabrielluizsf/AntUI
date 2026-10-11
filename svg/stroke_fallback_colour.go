package svg

import "github.com/gabrielluizsf/antui/canvas"

func (st Style) strokeFallbackColour(current canvas.Color) (canvas.Color, bool) {
	if !st.hasStrokeFallback {
		return 0, false
	}
	if st.strokeFallbackCurrent {
		return current, true
	}
	return st.strokeFallback, true
}
