package css

// parseGridAutoRepeatTracks reads the repeated list an auto-repeat introduces, without the repetition count.
func parseGridAutoRepeatTracks(raw string, ctx Units, depth int, fit bool) (GridTemplate, bool) {
	inner, ok := parseGridTemplateTracks(raw, ctx, depth+1)
	if !ok || len(inner) != 1 || inner[0].Auto {
		return nil, false
	}
	segment := inner[0]
	if len(segment.Before) > 0 || len(segment.After) > 0 {
		return nil, false
	}
	if _, ok := gridTrackMinimum(segment.Track, 0); !ok {
		return nil, false
	}
	segment.Auto, segment.Fit = true, fit
	return GridTemplate{segment}, true
}
