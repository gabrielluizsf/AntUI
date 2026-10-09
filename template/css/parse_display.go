package css

// parseDisplay reads a display keyword. A browser folds every table interior
// keyword into display:table because its anonymous box rules make the
// difference invisible; here the interior keywords are kept apart, so a
// stylesheet that styles a cell differently from the table around it is saying
// something the table layout can use. list-item, flow-root and the
// table-column boxes are accepted and read as the box they behave like.
func parseDisplay(raw string) (uint8, bool) {
	switch raw {
	case "block", "flow-root", "list-item":
		return DisplayBlock, true
	case "none":
		return DisplayNone, true
	case "inline":
		return DisplayInline, true
	case "inline-block", "inline-flex", "inline-grid":
		return DisplayInlineBlock, true
	case "flex":
		return DisplayFlex, true
	case "grid":
		return DisplayGrid, true
	case "table":
		return DisplayTable, true
	case "inline-table":
		return DisplayInlineTable, true
	case "table-row":
		return DisplayTableRow, true
	case "table-row-group":
		return DisplayTableRowGroup, true
	case "table-cell":
		return DisplayTableCell, true
	case "table-header-group":
		return DisplayTableHeaderGroup, true
	case "table-footer-group":
		return DisplayTableFooterGroup, true
	case "table-caption":
		return DisplayTableCaption, true
	case "table-column", "table-column-group":
		return DisplayTableColumn, true
	case "contents":
		return DisplayNone, true
	}
	return 0, false
}
