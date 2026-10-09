package css

// gridSegmentRepeat is how many tracks one segment stands for: a plain track is
// one, and an auto-repeat as many as fit the space the axis has.
func gridSegmentRepeat(segment GridSegment, available, gap, items int) int {
	if !segment.Auto {
		return 1
	}
	repeat := gridAutoRepeatCount(segment.Track, available, gap)
	if segment.Fit {
		repeat = min(repeat, max(items, 1))
	}
	return max(repeat, 1)
}
