package css

// parseGridTrack reads one track: a size, the auto keyword, or a line-name group.
func parseGridTrack(raw string, ctx Units) (GridTrack, bool) {
	name, inner, ok := gridFunction(raw)
	// fit-content() is a single track size that happens to be written as a
	// function, so it goes the way any other track size is read.
	if !ok || name == "fit-content" {
		size, good := parseGridTrackSize(raw, ctx, true)
		return GridTrack{Kind: GridTrackSingle, Size: size}, good
	}
	if name != "minmax" {
		return GridTrack{}, false
	}
	parts := splitGridList(inner, true)
	if len(parts) != 2 {
		return GridTrack{}, false
	}
	minSize, ok := parseGridTrackSize(parts[0], ctx, false)
	// The minimum of a minmax is a fixed breadth: auto, a length, or one of the
	// content keywords. fit-content() and a fraction are a track's maximum
	// only, so a minmax asking for one as its floor is not a minmax.
	if !ok || minSize.Kind == GridTrackFitContent {
		return GridTrack{}, false
	}
	maxSize, ok := parseGridTrackSize(parts[1], ctx, true)
	if !ok {
		return GridTrack{}, false
	}
	return GridTrack{Kind: GridTrackMinMax, Min: minSize, Max: maxSize}, true
}

// parseGridTemplate reads a grid-template value into its tracks. Track sizes and
// keywords fold to lowercase; line names keep the case the stylesheet wrote
// them with, the way CSS spells them.
