package css

// gridTrackMinimum is the smallest size a track can be forced down to, which is
// what the auto-repeat count is built from. Only a fixed minimum qualifies: a
// repeat with no definite minimum repeats once, as the grid grammar's
// auto-repeat requires.
func gridTrackMinimum(track GridTrack, available int) (int, bool) {
	size := track.Size
	if track.Kind == GridTrackMinMax {
		size = track.Min
	}
	if size.Kind != GridTrackLength {
		return 0, false
	}
	return max(size.Length.Resolve(Units{Width: available, Height: available}), 0), true
}
