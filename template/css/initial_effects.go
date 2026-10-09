package css

func initialEffects(st *Style, prop string) {
	switch prop {
	case "box-shadow":
		st.BoxShadow = nil
	case "text-shadow":
		st.TextShadow = nil
	case "outline":
		st.OutlineWidth = 0
		st.OutlineStyle = BorderNone
		st.OutlineColor = CurrentColor
	case "outline-width":
		st.OutlineWidth = 0
	case "outline-style":
		st.OutlineStyle = BorderNone
	case "outline-color":
		st.OutlineColor = CurrentColor
	case "outline-offset":
		st.OutlineOffset = 0
	case "filter":
		st.Filters = nil
	case "backdrop-filter":
		st.BackdropFilters = nil
	case "transform":
		st.Transform = nil
	case "transform-origin":
		st.TransformOrigin = InitialTransformOrigin
	}
}
