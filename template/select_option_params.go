package template

// SelectOptionParams holds parameters for painting a select option.
type SelectOptionParams struct {
	WidgetGeometry
	Label    string
	Selected bool
}
