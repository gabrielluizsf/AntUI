package css

import "strings"

// parseGridTemplateTracks reads the track list of grid-template, honouring nested line-name groups.
func parseGridTemplateTracks(raw string, ctx Units, depth int) (GridTemplate, bool) {
	if depth > 4 {
		return nil, false
	}
	parts := splitGridList(raw, false)
	if len(parts) == 0 {
		return nil, false
	}
	var out GridTemplate
	var before []string
	for _, part := range parts {
		if strings.HasPrefix(part, "[") {
			names, ok := parseGridLineNames(part)
			if !ok {
				return nil, false
			}
			if len(out) > 0 {
				last := &out[len(out)-1]
				last.After = append(last.After, names...)
				continue
			}
			before = append(before, names...)
			continue
		}
		segments, ok := parseGridTemplateSegment(part, ctx, depth)
		if !ok {
			return nil, false
		}
		segments[0].Before = append(before, segments[0].Before...)
		before = nil
		out = append(out, segments...)
		if len(out) > maxGridTracks {
			return nil, false
		}
	}
	if len(before) > 0 {
		return nil, false
	}
	return out, len(out) > 0
}

// parseGridTemplateSegment reads one item of a track list: a plain track, a
// repeat(n, …) expanded into n copies, or an auto-fill/auto-fit repeat left
// for the layout to count. The grid grammar allows an auto-repeat only once per
// axis, with a single track and a definite minimum, and it takes no line names.
