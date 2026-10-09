package css

// Cell reports whether the display value names a table cell, the box whose
// content a table row lines up in a column.
func (s Style) Cell() bool { return s.Display == DisplayTableCell }
