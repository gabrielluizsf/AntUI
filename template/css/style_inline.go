package css

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
