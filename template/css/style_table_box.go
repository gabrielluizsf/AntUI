package css

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
