package template

// tableChild is one thing a table draws between its rows: a row, or the
// caption the table paints above itself.
type tableChild struct {
	row     *tableRow
	caption *tableCell
}
