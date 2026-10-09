package css

// GridTrack is one column or row of the grid.
type GridTrackKind uint8

const (
	GridTrackSingle GridTrackKind = iota
	GridTrackMinMax
)

type GridTrack struct {
	Kind GridTrackKind
	Size GridTrackSize
	Min  GridTrackSize
	Max  GridTrackSize
}
