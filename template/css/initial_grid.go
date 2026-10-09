package css

func initialGrid(st *Style, prop string) {
	switch prop {
	case "grid-template-columns":
		st.GridTemplateColumns = nil
	case "grid-template-rows":
		st.GridTemplateRows = nil
	case "grid-template-areas":
		st.GridTemplateAreas = nil
	case "grid-auto-flow":
		st.GridAutoFlow = GridAutoFlowRow
		st.GridAutoFlowDense = false
	case "justify-items":
		st.JustifyItems = AlignStretch
	case "justify-self":
		st.JustifySelf = AlignAuto
	case "grid-column":
		st.GridColumn = GridPlacement{}
	case "grid-column-start":
		st.GridColumn.Start = GridLine{}
	case "grid-column-end":
		st.GridColumn.End = GridLine{}
	case "grid-row":
		st.GridRow = GridPlacement{}
	case "grid-row-start":
		st.GridRow.Start = GridLine{}
	case "grid-row-end":
		st.GridRow.End = GridLine{}
	case "grid-area":
		st.GridColumn = GridPlacement{}
		st.GridRow = GridPlacement{}
	}
}
