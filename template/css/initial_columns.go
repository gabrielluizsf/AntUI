package css

func initialColumns(st *Style, prop string) {
	switch prop {
	case "columns":
		st.ColumnCount = 0
		st.ColumnWidth = Auto()
	case "column-count":
		st.ColumnCount = 0
	case "column-width":
		st.ColumnWidth = Auto()
	case "column-fill":
		st.ColumnFill = ColumnFillBalance
	case "column-rule":
		st.ColumnRuleWidth = 0
		st.ColumnRuleStyle = BorderNone
		st.ColumnRuleColor = CurrentColor
	case "column-rule-width":
		st.ColumnRuleWidth = 0
	case "column-rule-style":
		st.ColumnRuleStyle = BorderNone
	case "column-rule-color":
		st.ColumnRuleColor = CurrentColor
	}
}
