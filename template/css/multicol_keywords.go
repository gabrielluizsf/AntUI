package css

// ColumnFill names how a multi-column container spends a height it was given.
// Balance, the initial value, spreads the content evenly so the columns come
// out the same height; auto fills the first column to the height before the
// next one starts.
const (
	ColumnFillBalance uint8 = iota
	ColumnFillAuto
)

// Break* names whether a box may be cut in the middle. The engine has one kind
// of fragmentation, the cut a multi-column container makes, and break-inside
// speaks to it: BreakAuto lets the flow cut a box wherever it runs out of
// room, BreakAvoid asks for the whole of it in one piece.
const (
	BreakAuto uint8 = iota
	BreakAvoid
)

// Float* names the side a box is taken out of the flow towards, and Clear* the
// sides a box asks to be pushed past. Both are read so the cascade has an
// answer for them; see doc.go for why this flow acts on neither.
const (
	FloatNone uint8 = iota
	FloatLeft
	FloatRight
	FloatInlineStart
	FloatInlineEnd
)

const (
	ClearNone uint8 = iota
	ClearLeft
	ClearRight
	ClearBoth
)
