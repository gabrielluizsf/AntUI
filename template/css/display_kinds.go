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
