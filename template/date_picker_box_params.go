package template

// DatePickerBoxParams holds parameters for painting a closed date picker.
type DatePickerBoxParams struct {
	WidgetGeometry
	Value string
	Open  bool
}