package css

import "github.com/gabrielluizsf/antui/canvas"

func initialSurface(st *Style, prop string) {
	switch prop {
	case "color":
		st.Color = ThemeInk
	case "opacity":
		st.Opacity = 1
	case "background":
		st.Background = canvas.Transparent
		st.BackgroundImages = nil
		st.BackgroundPos = nil
		st.BackgroundSize = nil
		st.BackgroundRepeat = nil
		st.BackgroundClip = nil
		st.BackgroundOrigin = nil
		st.BackgroundAttach = nil
	case "background-color":
		st.Background = canvas.Transparent
	case "background-image":
		st.BackgroundImages = nil
	case "background-position":
		st.BackgroundPos = nil
	case "background-size":
		st.BackgroundSize = nil
	case "background-repeat":
		st.BackgroundRepeat = nil
	case "background-clip":
		st.BackgroundClip = nil
	case "background-origin":
		st.BackgroundOrigin = nil
	case "background-attachment":
		st.BackgroundAttach = nil
	}
}
