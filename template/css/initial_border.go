package css

import "github.com/gabrielluizsf/antui/canvas"

func initialBorder(st *Style, prop string) {
	switch prop {
	case "border":
		st.BorderWidth = [4]int{}
		st.BoxStyle = [4]uint8{}
		st.BoxColor = [4]canvas.Color{}
	case "border-width":
		st.BorderWidth = [4]int{}
	case "border-style":
		st.BoxStyle = [4]uint8{}
	case "border-color":
		st.BoxColor = [4]canvas.Color{}
	case "border-radius":
		st.Radius = [4]int{}
		st.RadiusY = [4]int{}
	}
}
