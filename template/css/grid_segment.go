package css

// GridSegment is one track of a template, or one auto-repeat whose count the
// layout still owes.
type GridSegment struct {
	Track  GridTrack
	Auto   bool     // an auto-fill/auto-fit repeat: the count is the layout's
	Fit    bool     // auto-fit: repetitions no item lands in collapse
	Before []string // line names on the line before this track
	After  []string // line names on the line after this track
}
