package css

// Resolve turns the template into the explicit tracks of one axis. available is
// the container's content size along the axis and gap the grid gap there;
// items is how many in-flow items the grid holds, which is what auto-fit counts
// against so the repetitions nothing lands in are never created.
func (t GridTemplate) Resolve(available, gap, items int) GridAxis {
	tracks := make([]GridTrack, 0, len(t))
	names := make([][]string, 1, len(t)+1)
	for _, segment := range t {
		if len(names) > 0 {
			names[len(names)-1] = append(names[len(names)-1], segment.Before...)
		}
		for n := gridSegmentRepeat(segment, available, gap, items); n > 0; n-- {
			tracks = append(tracks, segment.Track)
			names = append(names, segment.After)
		}
	}
	return GridAxis{Tracks: tracks, Names: names}
}
