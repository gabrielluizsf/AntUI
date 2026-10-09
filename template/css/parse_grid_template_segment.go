package css

import (
	"strconv"
	"strings"
)

// parseGridTemplateSegment reads one track segment: a repeated group, an auto-repeat, or a plain track.
func parseGridTemplateSegment(raw string, ctx Units, depth int) (GridTemplate, bool) {
	name, inner, ok := gridFunction(raw)
	if !ok || name != "repeat" {
		track, good := parseGridTrack(raw, ctx)
		if !good {
			return nil, false
		}
		return GridTemplate{{Track: track}}, true
	}
	args := splitGridList(inner, true)
	if len(args) != 2 {
		return nil, false
	}
	if count, ok := parseGridAutoRepeat(args[0]); ok {
		return parseGridAutoRepeatTracks(args[1], ctx, depth, count)
	}
	count, err := strconv.Atoi(strings.TrimSpace(args[0]))
	if err != nil || count < 1 {
		return nil, false
	}
	expanded, ok := parseGridTemplateTracks(args[1], ctx, depth+1)
	if !ok || len(expanded) == 0 || count > maxGridTracks/len(expanded) {
		return nil, false
	}
	var out GridTemplate
	for n := 0; n < count; n++ {
		out = append(out, expanded...)
	}
	return out, true
}

// parseGridAutoRepeat recognises the auto-fill and auto-fit keywords, the only
// repeat() counts the layout works out instead of reading off a number.
