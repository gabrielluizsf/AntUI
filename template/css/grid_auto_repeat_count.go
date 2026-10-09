package css

// gridAutoRepeatCount is the repetition count of an auto-fill/auto-fit repeat:
// the largest number of the track's minimum that still fits the space available,
// gaps included, and never less than one.
func gridAutoRepeatCount(track GridTrack, available, gap int) int {
	minimum, ok := gridTrackMinimum(track, available)
	if !ok {
		return 1
	}
	step := minimum + gap
	if step <= 0 {
		return maxGridTracks
	}
	return min(max((available+gap)/step, 1), maxGridTracks)
}
