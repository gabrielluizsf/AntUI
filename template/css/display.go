package css

// Display is the CSS display type of a box: what kind of box it is, and for
// the container types, what lays the widgets drawn inside it out. The value is
// both the outer type (block, inline, flex, grid, table) and the inner one for
// the boxes that only exist inside a table (a row, a cell, a caption), which
// is what lets a stylesheet say "this widget is a cell" and have the table
// layout hear it.
const (
	DisplayBlock uint8 = iota
	DisplayNone
	DisplayInline
	DisplayInlineBlock
	DisplayFlex
	DisplayGrid
	DisplayTable
	DisplayInlineTable
	DisplayTableRow
	DisplayTableRowGroup
	DisplayTableCell
	DisplayTableHeaderGroup
	DisplayTableFooterGroup
	DisplayTableCaption
	DisplayTableColumn
	DisplayTableColumnGroup
)

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

// Inline reports whether the box flows along a line with its neighbours. An
// inline-table is a table that stays on the line, as inline-block is a block
// that does.
func (s Style) Inline() bool {
	switch s.Display {
	case DisplayInline, DisplayInlineBlock, DisplayInlineTable:
		return true
	}
	return false
}

// TableBox reports whether the display value names a box that only exists
// inside a table: a row or a row group. The table layout uses it to recognise
// a row among its children, the way CSS does before it wraps anything.
func (s Style) TableBox() bool {
	switch s.Display {
	case DisplayTableRow, DisplayTableRowGroup, DisplayTableHeaderGroup, DisplayTableFooterGroup:
		return true
	}
	return false
}

// Cell reports whether the display value names a table cell, the box whose
// content a table row lines up in a column.
func (s Style) Cell() bool { return s.Display == DisplayTableCell }

// Caption reports whether the display value names a table caption, the box a
// table paints above itself, across its whole width.
func (s Style) Caption() bool { return s.Display == DisplayTableCaption }
