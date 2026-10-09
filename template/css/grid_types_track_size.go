package css

// GridTrackSize is one track length: a range of sizes
// (the min is the max when unsplit), the auto keyword, or a flexible "fr" weight.
type GridTrackSizeKind uint8

const (
	GridTrackAuto GridTrackSizeKind = iota
	GridTrackLength
	GridTrackFlexible
	GridTrackMinContent
	GridTrackMaxContent
	GridTrackFitContent
)

// GridTrackSize is a <track-size>: a length, a flexible fraction, the auto
// keyword, or one of the content keywords — min-content, max-content and
// fit-content(), whose limit travels in Length.
type GridTrackSize struct {
	Kind   GridTrackSizeKind
	Length Length
	Fr     float64
}
