package css

func applyAlignDecl(e *declEnv) []string {
	st, set, raw, ctx, note := e.st, e.set, e.raw, e.ctx, e.note
	switch e.prop {
	case "justify-content":
		if v, ok := parseJustifyContent(raw); ok {
			st.JustifyContent = v
			if grid, good := parseGridJustifyContent(raw); good {
				st.GridJustifyContent = grid
			}
			note()
		}
	case "align-items":
		if v, ok := parseAlignItems(raw); ok {
			st.AlignItems = v
			note()
		}
	case "align-self":
		if v, ok := parseAlignSelf(raw); ok {
			st.AlignSelf = v
			note()
		}
	case "align-content":
		if v, ok := parseAlignContent(raw); ok {
			st.AlignContent = v
			note()
		}
	case "gap":
		if r, c, ok := parseGap(raw, *ctx); ok {
			st.RowGap, st.ColumnGap = r, c
			note()
			set["row-gap"] = true
			set["column-gap"] = true
		}
	case "row-gap":
		if r, _, ok := parseGap(raw, *ctx); ok {
			st.RowGap = r
			note()
			set["gap"] = true
		}
	case "column-gap":
		if _, c, ok := parseGap(raw, *ctx); ok {
			st.ColumnGap = c
			note()
			set["gap"] = true
		}
	}
	return nil
}
