package css

// inheritOne copies one computed value from the body style into the element
// and records that the property came by inheritance. The body holds the
// winning cascade value already—including a value the body itself inherited
// from nothing, which is the theme's zero.
func inheritOne(st *Style, set map[string]bool, body Style, prop string) {
	switch prop {
	case "color":
		st.Color = body.Color
	case "font-size":
		st.FontSize = body.FontSize
	case "font-weight":
		st.FontWeight = body.FontWeight
	case "font-style":
		st.FontStyle = body.FontStyle
	case "font-family":
		st.FontFamily = body.FontFamily
	case "line-height":
		st.LineHeight = body.LineHeight
	case "letter-spacing":
		st.LetterSpacing = body.LetterSpacing
	case "word-spacing":
		st.WordSpacing = body.WordSpacing
	case "text-align":
		st.TextAlign = body.TextAlign
	case "text-transform":
		st.TextTransform = body.TextTransform
	case "white-space":
		st.WhiteSpace = body.WhiteSpace
	case "overflow-wrap":
		st.OverflowWrap = body.OverflowWrap
	case "visibility":
		st.Visibility = body.Visibility
	case "cursor":
		st.Cursor = body.Cursor
	case "text-shadow":
		st.TextShadow = body.TextShadow
	}
	set[prop] = true
	delete(st.inherit, prop)
}
