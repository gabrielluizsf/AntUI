package template

// SelectParams holds parameters for painting a select dropdown.
type SelectParams struct {
	WidgetGeometry
	Value string
	Open  bool
}
