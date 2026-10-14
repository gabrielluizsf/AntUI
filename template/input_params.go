package template

// InputParams holds parameters for painting a single-line text input.
type InputParams struct {
	WidgetGeometry
	Text   string
	Cursor int
	Caret  bool
}
