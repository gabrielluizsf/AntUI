package css

func initialFlex(st *Style, prop string) {
	switch prop {
	case "flex-direction":
		st.FlexDirection = FlexDirectionRow
	case "flex-wrap":
		st.FlexWrap = FlexWrapNowrap
	case "flex-flow":
		st.FlexDirection = FlexDirectionRow
		st.FlexWrap = FlexWrapNowrap
	case "justify-content":
		st.JustifyContent = JustifyFlexStart
		st.GridJustifyContent = ContentStretch
	case "align-items":
		st.AlignItems = AlignStretch
	case "align-self":
		st.AlignSelf = AlignAuto
	case "align-content":
		st.AlignContent = ContentStretch
	case "gap":
		st.RowGap = Zero()
		st.ColumnGap = Zero()
	case "row-gap":
		st.RowGap = Zero()
	case "column-gap":
		st.ColumnGap = Zero()
	case "order":
		st.Order = 0
	case "flex":
		st.FlexGrow = 0
		st.FlexShrink = 1
		st.FlexBasis = Auto()
	case "flex-grow":
		st.FlexGrow = 0
	case "flex-shrink":
		st.FlexShrink = 1
	case "flex-basis":
		st.FlexBasis = Auto()
	}
}
