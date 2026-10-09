package css

func initialLayout(st *Style, prop string) {
	switch prop {
	case "display":
		st.Display = DisplayBlock
	case "position":
		st.Position = PositionStatic
	case "top":
		st.Top = Auto()
	case "right":
		st.Right = Auto()
	case "bottom":
		st.Bottom = Auto()
	case "left":
		st.Left = Auto()
	case "z-index":
		st.ZIndex = 0
	case "overflow":
		st.Overflow = [2]uint8{OverflowVisible, OverflowVisible}
	case "overflow-x":
		st.Overflow[0] = OverflowVisible
	case "overflow-y":
		st.Overflow[1] = OverflowVisible
	case "visibility":
		st.Visibility = VisibilityVisible
	case "pointer-events":
		st.PointerEvents = PointerEventsAuto
	case "cursor":
		st.Cursor = CursorDefault
	case "break-inside", "page-break-inside":
		st.BreakInside = BreakAuto
	case "float":
		st.Float = FloatNone
	case "clear":
		st.Clear = ClearNone
	}
}
