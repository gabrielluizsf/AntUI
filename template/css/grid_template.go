package css

// The template a grid axis is built from: grid-template-columns or
// grid-template-rows read into its tracks, and the line names that stand
// between them.

import "strings"

// maxGridTracks caps how many tracks one axis may end up with, so a
// pathological repeat() cannot make the solver walk a million entries.
const maxGridTracks = 1000

// GridTemplate is a parsed grid-template value: one segment per track, in
// declaration order, each carrying the line names standing on the lines to its
// left and right. repeat(n, …) is already expanded here, where the count is
// written down; repeat(auto-fill | auto-fit, …) stays a single segment and
// [GridTemplate.Resolve] decides how many times it repeats.
type GridTemplate []GridSegment

// GridSegment is one track of a template, or one auto-repeat whose count the
// layout still owes.
type GridSegment struct {
	Track  GridTrack
	Auto   bool     // an auto-fill/auto-fit repeat: the count is the layout's
	Fit    bool     // auto-fit: repetitions no item lands in collapse
	Before []string // line names on the line before this track
	After  []string // line names on the line after this track
}

// GridAxis is one axis of a grid with its template resolved: the explicit
// tracks, and the names standing on the lines between them — one entry per
// line, the first being the line before the first track.
type GridAxis struct {
	Tracks []GridTrack
	Names  [][]string
}

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

// LineIndex finds the line carrying a name: the first one counting from the
// start, the last one counting from the end — the way a placement resolves its
// start line and its end line. Implicit lines hold no names, so a name no
// explicit line carries is simply not found.
func (a GridAxis) LineIndex(name string, fromEnd bool) (int, bool) {
	if fromEnd {
		for i := len(a.Names) - 1; i >= 0; i-- {
			if gridHasName(a.Names[i], name) {
				return i, true
			}
		}
		return 0, false
	}
	for i, names := range a.Names {
		if gridHasName(names, name) {
			return i, true
		}
	}
	return 0, false
}

// gridHasName reports whether one line carries a name.
func gridHasName(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}

// parseGridLineNames reads a "[…]" line-name list. Names are case-sensitive,
// the way CSS spells them, and the keywords the grid grammar already spent are
// refused: a line named "auto" or "span" would read as a placement, not a name.
func parseGridLineNames(raw string) ([]string, bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) < 3 || raw[0] != '[' || raw[len(raw)-1] != ']' {
		return nil, false
	}
	fields := strings.Fields(raw[1 : len(raw)-1])
	if len(fields) == 0 {
		return nil, false
	}
	for _, field := range fields {
		if !validGridName(field) {
			return nil, false
		}
	}
	return fields, true
}

// validGridName is a line name: an identifier that is none of the keywords the
// grid grammar spells elsewhere.
func validGridName(raw string) bool {
	if !validGridIdent(raw) {
		return false
	}
	switch strings.ToLower(raw) {
	case "auto", "span", "min-content", "max-content", "fit-content", "minmax",
		"repeat", "subgrid", "initial", "inherit", "unset", "none":
		return false
	}
	return true
}
