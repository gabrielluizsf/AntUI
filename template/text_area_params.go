package template

// TextAreaParams holds parameters for painting a multi-line text area.
type TextAreaParams struct {
	WidgetGeometry
	Text   string
	Cursor int
	Caret  bool
}
