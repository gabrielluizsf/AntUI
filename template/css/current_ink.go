package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func currentInk(st *Style) canvas.Color {
	if st.Set["color"] && st.Color != CurrentColor {
		return st.Color
	}
	return 0
}
