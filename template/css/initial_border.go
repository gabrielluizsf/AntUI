package css

import "github.com/gabrielluizsf/antui/canvas"

func initialBorder(st *Style, prop string) {
	switch prop {
	case "border":
		st.BorderWidth = [sideCount]int{}
		st.BoxStyle = [sideCount]uint8{}
		st.BoxColor = [sideCount]canvas.Color{}
	case "border-width":
		st.BorderWidth = [sideCount]int{}
	case "border-style":
		st.BoxStyle = [sideCount]uint8{}
	case "border-color":
		st.BoxColor = [sideCount]canvas.Color{}
	case "border-radius":
		st.Radius = [sideCount]int{}
		st.RadiusY = [sideCount]int{}
	}
}
