package css

// Caption reports whether the display value names a table caption, the box a
// table paints above itself, across its whole width.
func (s Style) Caption() bool { return s.Display == DisplayTableCaption }
