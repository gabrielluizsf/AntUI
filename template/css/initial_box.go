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
		st.Margin = [sideCount]Length{Zero(), Zero(), Zero(), Zero()}
	case "margin-top":
		st.Margin[SideTop] = Zero()
	case "margin-right":
		st.Margin[SideRight] = Zero()
	case "margin-bottom":
		st.Margin[SideBottom] = Zero()
	case "margin-left":
		st.Margin[SideLeft] = Zero()
	case "padding":
		st.Padding = [sideCount]Length{Zero(), Zero(), Zero(), Zero()}
	case "padding-top":
		st.Padding[SideTop] = Zero()
	case "padding-right":
		st.Padding[SideRight] = Zero()
	case "padding-bottom":
		st.Padding[SideBottom] = Zero()
	case "padding-left":
		st.Padding[SideLeft] = Zero()
	case "box-sizing":
		st.BoxSizing = false
	}
}
