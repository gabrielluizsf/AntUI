package css

import "github.com/gabrielluizsf/antui/canvas"

// columnStyle carries the multi-column and flow-control fields of a Style.
// ColumnCount is how many columns the content is split into, ColumnWidth how
// wide each one is when the count is left to the width, ColumnFill how a
// declared height is spent. The rule in the gutter is ColumnRuleWidth,
// ColumnRuleStyle (the Border* constants) and ColumnRuleColor. The gutter
// itself is ColumnGap, the same property flex and grid read as their gap.
//
// BreakInside is whether a box may be cut in the middle where it is laid
// out: in a column, in a page or along a line. BreakAuto lets the flow cut
// it wherever it runs out of room, BreakAvoid asks for the whole box in
// one piece.
//
// Float is the side a box is taken out of the flow towards and Clear the
// sides a box asks to be pushed past. Both are read, so a stylesheet that
// says float: right is heard and not warned about as a misspelling; see
// doc.go for why this flow has nothing to float inside.
type columnStyle struct {
	ColumnCount     int
	ColumnWidth     Length
	ColumnFill      uint8 // one of the ColumnFill* constants
	ColumnRuleWidth int
	ColumnRuleStyle uint8
	ColumnRuleColor canvas.Color

	BreakInside uint8 // one of the Break* constants

	Float uint8 // one of the Float* constants
	Clear uint8 // one of the Clear* constants
}
