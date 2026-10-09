package css

// GridLine is one boundary of the rail grid: auto, an index or a named line (a leading "-" counts from the end).
type GridLineKind uint8

const (
	GridLineAuto GridLineKind = iota
	GridLineIndex
	GridLineSpan
	GridLineName
)

type GridLine struct {
	Kind GridLineKind
	// Index is the line number a placement names, positive or negative, or the
	// count a bare span stands for.
	Index int
	// Name is the line name, when the placement gave one.
	Name string
	// Span is a count that follows the line, as in "2 span 3" or
	// "sidebar-start span 2"; zero when the placement gave none.
	Span int
	// Backward marks a name counted from the end, written "-name".
	Backward bool
}
