package svg

import "github.com/gabrielluizsf/antui/canvas"

// fillFallbackColour and strokeFallbackColour are the colour after a `url(...)`
// where the reference is not one the drawing can follow — never there, or there
// and with nothing to paint with — which is the second choice the file gave, if
// it gave one. Writing it `currentColor` keeps the keyword: the colour of the
// widget is only known when it is painted.
func (st Style) fillFallbackColour(current canvas.Color) (canvas.Color, bool) {
	if !st.hasFillFallback {
		return 0, false
	}
	if st.fillFallbackCurrent {
		return current, true
	}
	return st.fillFallback, true
}
