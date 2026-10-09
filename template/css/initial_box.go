package css

func initialBox(st *Style, prop string) {
	switch prop {
	case "width":
		st.Width = Auto()
	case "min-width":
		st.MinWidth = Auto()
	case "height":
		st.Height = Auto()
	case "min-height":
		st.MinHeight = Auto()
	case "max-width":
		st.MaxWidth = Length{u: unitNone}
	case "max-height":
		st.MaxHeight = Length{u: unitNone}
	case "margin":
		st.Margin = [4]Length{Zero(), Zero(), Zero(), Zero()}
	case "margin-top":
		st.Margin[0] = Zero()
	case "margin-right":
		st.Margin[1] = Zero()
	case "margin-bottom":
		st.Margin[2] = Zero()
	case "margin-left":
		st.Margin[3] = Zero()
	case "padding":
		st.Padding = [4]Length{Zero(), Zero(), Zero(), Zero()}
	case "padding-top":
		st.Padding[0] = Zero()
	case "padding-right":
		st.Padding[1] = Zero()
	case "padding-bottom":
		st.Padding[2] = Zero()
	case "padding-left":
		st.Padding[3] = Zero()
	case "box-sizing":
		st.BoxSizing = false
	}
}
