package css

func applyGridPlacementDecl(e *declEnv) []string {
	st, set, raw, prop, note := e.st, e.set, e.raw, e.prop, e.note
	switch prop {
	case "grid-column":
		if v, ok := parseGridPlacement(raw); ok {
			st.GridColumn = v
			note()
			set["grid-column-start"] = true
			set["grid-column-end"] = true
		}
	case "grid-row":
		if v, ok := parseGridPlacement(raw); ok {
			st.GridRow = v
			note()
			set["grid-row-start"] = true
			set["grid-row-end"] = true
		}
	case "grid-column-start", "grid-column-end":
		if v, ok := parseGridLine(raw); ok {
			if prop == "grid-column-start" {
				st.GridColumn.Start = v
			} else {
				st.GridColumn.End = v
			}
			note()
			set["grid-column"] = true
		}
	case "grid-row-start", "grid-row-end":
		if v, ok := parseGridLine(raw); ok {
			if prop == "grid-row-start" {
				st.GridRow.Start = v
			} else {
				st.GridRow.End = v
			}
			note()
			set["grid-row"] = true
		}
	case "grid-area":
		if column, row, ok := parseGridArea(raw); ok {
			st.GridColumn = column
			st.GridRow = row
			note()
			set["grid-column"] = true
			set["grid-row"] = true
			set["grid-column-start"] = true
			set["grid-column-end"] = true
			set["grid-row-start"] = true
			set["grid-row-end"] = true
		}
	}
	return nil
}
