package css

// GridPlacement is one element position: a start line and an optional end line, split by a slash.
type GridPlacement struct {
	Start GridLine
	End   GridLine
}

const (
	GridAutoFlowRow uint8 = iota
	GridAutoFlowColumn
)
