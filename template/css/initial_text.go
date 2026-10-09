package css

func initialText(st *Style, prop string) {
	switch prop {
	case "font-size":
		st.FontSize = 0
	case "text-align":
		st.TextAlign = 0
	case "font-weight":
		st.FontWeight = 0
	case "font-style":
		st.FontStyle = FontStyleNormal
	case "font-family":
		st.FontFamily = ""
	case "line-height":
		st.LineHeight = 0
	case "letter-spacing":
		st.LetterSpacing = Zero()
	case "word-spacing":
		st.WordSpacing = Zero()
	case "text-transform":
		st.TextTransform = TextTransformNone
	case "text-decoration", "text-decoration-line":
		st.TextDecoration = 0
	case "white-space":
		st.WhiteSpace = WhiteSpaceNormal
	case "overflow-wrap", "word-wrap":
		st.OverflowWrap = OverflowWrapNormal
	case "text-overflow":
		st.TextOverflow = TextOverflowClip
	case "vertical-align":
		st.VerticalAlign = VerticalAlignBaseline
		st.BaselineShift = Length{}
	}
}
