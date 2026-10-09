package css

// GridAxis is one axis of a grid with its template resolved: the explicit
// tracks, and the names standing on the lines between them — one entry per
// line, the first being the line before the first track.
type GridAxis struct {
	Tracks []GridTrack
	Names  [][]string
}
