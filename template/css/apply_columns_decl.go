package css

func applyColumnsDecl(e *declEnv) []string {
	st, set, raw, ctx, note := e.st, e.set, e.raw, e.ctx, e.note
	switch e.prop {
	case "columns":
		if n, l, ok := parseColumns(raw, *ctx); ok {
			st.ColumnCount, st.ColumnWidth = n, l
			note()
			set["column-count"] = true
			set["column-width"] = true
		}
	case "column-count":
		if v, ok := parseColumnCount(raw); ok {
			st.ColumnCount = v
			note()
		}
	case "column-width":
		if v, ok := parseColumnWidth(raw, *ctx); ok {
			st.ColumnWidth = v
			note()
		}
	case "column-fill":
		if v, ok := parseColumnFill(raw); ok {
			st.ColumnFill = v
			note()
		}
	case "column-rule":
		applyColumnRule(st, set, raw, note)
	case "column-rule-width":
		if l, err := parseLength(raw); err == nil && l.u == unitPx {
			st.ColumnRuleWidth = max(l.Px(0), 0)
			note()
		}
	case "column-rule-style":
		if v, ok := parseBorderStyle(raw); ok {
			st.ColumnRuleStyle = v
			note()
		}
	case "column-rule-color":
		if c, err := ParseColor(raw); err == nil {
			st.ColumnRuleColor = c
			note()
		}
	}
	return nil
}
