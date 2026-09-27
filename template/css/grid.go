package css

import (
	"math"
	"strconv"
	"strings"
)

type GridTrackSizeKind uint8

const (
	GridTrackAuto GridTrackSizeKind = iota
	GridTrackLength
	GridTrackFlexible
	GridTrackMinContent
	GridTrackMaxContent
	GridTrackFitContent
)

// GridTrackSize is a <track-size>: a length, a flexible fraction, the auto
// keyword, or one of the content keywords — min-content, max-content and
// fit-content(), whose limit travels in Length.
type GridTrackSize struct {
	Kind   GridTrackSizeKind
	Length Length
	Fr     float64
}

type GridTrackKind uint8

const (
	GridTrackSingle GridTrackKind = iota
	GridTrackMinMax
)

type GridTrack struct {
	Kind GridTrackKind
	Size GridTrackSize
	Min  GridTrackSize
	Max  GridTrackSize
}

type GridLineKind uint8

const (
	GridLineAuto GridLineKind = iota
	GridLineIndex
	GridLineSpan
	GridLineName
)

type GridLine struct {
	Kind GridLineKind
	// Index is the line number a placement names, positive or negative, or the
	// count a bare span stands for.
	Index int
	// Name is the line name, when the placement gave one.
	Name string
	// Span is a count that follows the line, as in "2 span 3" or
	// "sidebar-start span 2"; zero when the placement gave none.
	Span int
	// Backward marks a name counted from the end, written "-name".
	Backward bool
}

type GridPlacement struct {
	Start GridLine
	End   GridLine
}

const (
	GridAutoFlowRow uint8 = iota
	GridAutoFlowColumn
)

func parseGridTrackSize(raw string, ctx Units, flexible bool) (GridTrackSize, bool) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	switch raw {
	case "auto":
		return GridTrackSize{Kind: GridTrackAuto}, true
	case "min-content":
		return GridTrackSize{Kind: GridTrackMinContent}, true
	case "max-content":
		return GridTrackSize{Kind: GridTrackMaxContent}, true
	}
	if name, inner, ok := gridFunction(raw); ok && name == "fit-content" {
		limit, err := parseLengthAt(inner, ctx)
		if err != nil || limit.Auto() || limit.None() || limit.value < 0 {
			return GridTrackSize{}, false
		}
		return GridTrackSize{Kind: GridTrackFitContent, Length: limit}, true
	}
	if flexible && strings.HasSuffix(raw, "fr") {
		v, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(raw, "fr")), 64)
		if err != nil || v <= 0 || math.IsNaN(v) || math.IsInf(v, 0) {
			return GridTrackSize{}, false
		}
		return GridTrackSize{Kind: GridTrackFlexible, Fr: v}, true
	}
	l, err := parseLengthAt(raw, ctx)
	if err != nil || l.Auto() || l.None() || l.value < 0 {
		return GridTrackSize{}, false
	}
	return GridTrackSize{Kind: GridTrackLength, Length: l}, true
}

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
func parseGridTemplate(raw string, ctx Units) (GridTemplate, bool) {
	raw = strings.TrimSpace(raw)
	if strings.EqualFold(raw, "none") {
		return nil, true
	}
	return parseGridTemplateTracks(raw, ctx, 0)
}

// parseGridTemplateTracks reads one track list — the value itself, or the body
// of a repeat() — into the segments it stands for. A "[…]" group names the
// line to its left: names seen before any track wait for it, and names seen
// after one close the line that track ends.
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
func parseGridAutoRepeat(raw string) (fit bool, ok bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "auto-fill":
		return false, true
	case "auto-fit":
		return true, true
	}
	return false, false
}

// parseGridAutoRepeatTracks reads the track an auto-repeat stands for: one
// track, no line names, and a minimum small enough to count against the space
// the axis has.
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

func parseGridTemplateAreas(raw string) ([][]string, bool) {
	raw = strings.TrimSpace(raw)
	if strings.EqualFold(raw, "none") {
		return nil, true
	}
	rows := splitGridList(raw, false)
	if len(rows) == 0 {
		return nil, false
	}
	out := make([][]string, 0, len(rows))
	width := -1
	for _, row := range rows {
		if len(row) < 2 || (row[0] != '"' && row[0] != '\'') || row[len(row)-1] != row[0] || !closedGridString(row) {
			return nil, false
		}
		cells := splitGridList(row[1:len(row)-1], false)
		if len(cells) == 0 {
			return nil, false
		}
		if width < 0 {
			width = len(cells)
		} else if len(cells) != width {
			return nil, false
		}
		for _, cell := range cells {
			if cell != "." && cell != "..." && !validGridIdent(cell) {
				return nil, false
			}
		}
		out = append(out, cells)
	}
	bounds := make(map[string][4]int)
	for row, cells := range out {
		for col, name := range cells {
			if name == "." || name == "..." {
				continue
			}
			b, exists := bounds[name]
			if !exists {
				bounds[name] = [4]int{col, row, col, row}
				continue
			}
			b[0] = min(b[0], col)
			b[1] = min(b[1], row)
			b[2] = max(b[2], col)
			b[3] = max(b[3], row)
			bounds[name] = b
		}
	}
	for name, b := range bounds {
		for row := b[1]; row <= b[3]; row++ {
			for col := b[0]; col <= b[2]; col++ {
				if out[row][col] != name {
					return nil, false
				}
			}
		}
	}
	return out, true
}

func parseGridAutoFlow(raw string) (uint8, bool, bool) {
	words := strings.Fields(raw)
	if len(words) == 0 || len(words) > 2 {
		return 0, false, false
	}
	var flow uint8
	var dense, haveFlow, haveDense bool
	for _, word := range words {
		switch word {
		case "row":
			if haveFlow {
				return 0, false, false
			}
			flow, haveFlow = GridAutoFlowRow, true
		case "column":
			if haveFlow {
				return 0, false, false
			}
			flow, haveFlow = GridAutoFlowColumn, true
		case "dense":
			if haveDense {
				return 0, false, false
			}
			dense, haveDense = true, true
		default:
			return 0, false, false
		}
	}
	return flow, dense, haveFlow
}

// parseGridLine reads one <grid-line>: "auto", a line number counting from the
// start or the end, a line name, or any of them followed by a span.
func parseGridLine(raw string) (GridLine, bool) {
	raw = strings.TrimSpace(raw)
	if strings.EqualFold(raw, "auto") {
		return GridLine{Kind: GridLineAuto}, true
	}
	if n, err := strconv.Atoi(raw); err == nil {
		if n == 0 {
			return GridLine{}, false
		}
		return GridLine{Kind: GridLineIndex, Index: n}, true
	}
	return parseGridLineName(raw)
}

// parseGridLineName reads a line name and the span that may follow it, in any of
// the orders CSS allows: "sidebar-start", "-sidebar-start", "span 2", "2 span",
// "2 span 3", "sidebar-start span 2". The name keeps the case it was written
// with; only the span keyword folds.
func parseGridLineName(raw string) (GridLine, bool) {
	fields := strings.Fields(raw)
	if len(fields) == 0 || len(fields) > 3 {
		return GridLine{}, false
	}
	// The span keyword tells a number from a line to a count: a number written
	// before it names the line, a number written after it counts the tracks
	// spanned. A name is always the line, and a leading "-" counts that name
	// from the end of the axis.
	line := GridLine{Kind: GridLineAuto}
	spans, counts, count := 0, 0, 0
	for _, field := range fields {
		if strings.EqualFold(field, "span") {
			spans++
			continue
		}
		if n, err := strconv.Atoi(field); err == nil {
			if spans == 0 {
				if n == 0 || line.Kind != GridLineAuto {
					return GridLine{}, false
				}
				line = GridLine{Kind: GridLineIndex, Index: n}
				continue
			}
			if counts > 0 || n < 1 {
				return GridLine{}, false
			}
			counts, count = 1, n
			continue
		}
		if line.Kind != GridLineAuto {
			return GridLine{}, false
		}
		name, backward := gridBackwardName(field)
		if !validGridName(name) {
			return GridLine{}, false
		}
		line = GridLine{Kind: GridLineName, Name: name, Backward: backward}
	}
	if spans > 1 {
		return GridLine{}, false
	}
	if spans == 0 {
		// A plain line, no span: only a name reaches here, a lone number
		// having been read by parseGridLine itself.
		return line, line.Kind == GridLineName
	}
	if count == 0 {
		count = 1
	}
	if line.Kind == GridLineAuto {
		return GridLine{Kind: GridLineSpan, Index: count}, true
	}
	line.Span = count
	return line, true
}

// gridBackwardName splits a leading "-" off a line name, the way CSS counts a
// named line from the end of the axis.
func gridBackwardName(raw string) (string, bool) {
	if strings.HasPrefix(raw, "-") {
		return raw[1:], true
	}
	return raw, false
}

func parseGridPlacement(raw string) (GridPlacement, bool) {
	parts := splitGridSlashes(raw)
	if len(parts) == 0 || len(parts) > 2 {
		return GridPlacement{}, false
	}
	placement := GridPlacement{}
	var ok bool
	if placement.Start, ok = parseGridLine(parts[0]); !ok {
		return GridPlacement{}, false
	}
	if len(parts) == 2 {
		if placement.End, ok = parseGridLine(parts[1]); !ok {
			return GridPlacement{}, false
		}
	}
	return placement, true
}

// parseGridArea reads a grid-area: a single value is one line — a name standing
// for the item's four lines, or a row line; two to four values are line
// placements written row first — "row-start / column-start [ / row-end /
// column-end ]".
func parseGridArea(raw string) (GridPlacement, GridPlacement, bool) {
	raw = strings.TrimSpace(raw)
	if strings.EqualFold(raw, "none") || strings.EqualFold(raw, "auto") {
		return GridPlacement{}, GridPlacement{}, true
	}
	parts := splitGridSlashes(raw)
	if len(parts) == 0 || len(parts) > 4 {
		return GridPlacement{}, GridPlacement{}, false
	}
	lines := make([]GridLine, len(parts))
	for i, part := range parts {
		line, ok := parseGridLine(part)
		if !ok {
			return GridPlacement{}, GridPlacement{}, false
		}
		lines[i] = line
	}
	if len(lines) == 1 {
		area := GridPlacement{Start: lines[0]}
		if lines[0].Kind == GridLineName {
			// A lone name is the name of a template area: it stands for all
			// four lines at once.
			area.End = lines[0]
			return area, area, true
		}
		return GridPlacement{}, area, true
	}
	row := GridPlacement{Start: lines[0]}
	column := GridPlacement{Start: lines[1]}
	if len(lines) >= 3 {
		row.End = lines[2]
	}
	if len(lines) == 4 {
		column.End = lines[3]
	}
	return column, row, true
}

func parseGridJustifyContent(raw string) (uint8, bool) {
	switch raw {
	case "normal", "stretch", "left", "start", "flex-start":
		if raw == "normal" || raw == "stretch" {
			return ContentStretch, true
		}
		return ContentFlexStart, true
	case "right", "end", "flex-end":
		return ContentFlexEnd, true
	case "center":
		return ContentCenter, true
	case "space-between":
		return ContentSpaceBetween, true
	case "space-around":
		return ContentSpaceAround, true
	case "space-evenly":
		return ContentSpaceEvenly, true
	}
	return 0, false
}

func parseJustifyItems(raw string) (uint8, bool) {
	switch raw {
	case "normal", "stretch", "left":
		if raw == "normal" || raw == "stretch" {
			return AlignStretch, true
		}
		return AlignFlexStart, true
	case "start", "self-start":
		return AlignFlexStart, true
	case "right", "end", "self-end":
		return AlignFlexEnd, true
	case "center":
		return AlignCenter, true
	}
	return 0, false
}

func parseJustifySelf(raw string) (uint8, bool) {
	if raw == "auto" {
		return AlignAuto, true
	}
	return parseJustifyItems(raw)
}

// gridFunction splits a track function call — "minmax(0, 1fr)", "repeat(2, …)",
// "fit-content(200px)" — into its lowercased name and its arguments. A function
// name is case-insensitive, so it is folded for the caller to match.
func gridFunction(raw string) (string, string, bool) {
	raw = strings.TrimSpace(raw)
	open := strings.IndexByte(raw, '(')
	if open <= 0 || !strings.HasSuffix(raw, ")") || parenEnd(raw, open) != len(raw)-1 {
		return "", "", false
	}
	name := strings.ToLower(strings.TrimSpace(raw[:open]))
	return name, strings.TrimSpace(raw[open+1 : len(raw)-1]), true
}

func splitGridList(raw string, comma bool) []string {
	var out []string
	start, depth := -1, 0
	var quote byte
	flush := func(end int) {
		if start >= 0 {
			out = append(out, strings.TrimSpace(raw[start:end]))
			start = -1
			return
		}
		if comma {
			out = append(out, "")
		}
	}
	for i := 0; i < len(raw); i++ {
		ch := raw[i]
		if quote != 0 {
			if ch == '\\' {
				i++
			} else if ch == quote {
				quote = 0
			}
			continue
		}
		if ch == '"' || ch == '\'' {
			if start < 0 {
				start = i
			}
			quote = ch
			continue
		}
		switch ch {
		case '(':
			if start < 0 {
				start = i
			}
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',', ' ', '\t', '\n', '\r':
			if depth == 0 && (comma && ch == ',' || !comma && ch != ',') {
				flush(i)
			}
		default:
			if start < 0 {
				start = i
			}
		}
	}
	flush(len(raw))
	return out
}

func splitGridSlashes(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	start, depth := -1, 0
	hadSlash := false
	for i := 0; i < len(raw); i++ {
		switch raw[i] {
		case '(':
			if start < 0 {
				start = i
			}
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case '/':
			if depth == 0 {
				if start < 0 {
					start = i
				}
				out = append(out, strings.TrimSpace(raw[start:i]))
				start = -1
				hadSlash = true
			}
		default:
			if start < 0 {
				start = i
			}
		}
	}
	if start >= 0 {
		out = append(out, strings.TrimSpace(raw[start:]))
	} else if hadSlash {
		out = append(out, "")
	}
	return out
}

func closedGridString(raw string) bool {
	if len(raw) < 2 || (raw[0] != '"' && raw[0] != '\'') {
		return false
	}
	quote := raw[0]
	for i := 1; i < len(raw)-1; i++ {
		if raw[i] == '\\' {
			i++
			continue
		}
		if raw[i] == quote {
			return false
		}
	}
	return true
}

func validGridIdent(raw string) bool {
	if raw == "" {
		return false
	}
	for i := 0; i < len(raw); i++ {
		ch := raw[i]
		if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '-' || ch == '_' || ch >= 0x80 {
			continue
		}
		if ch == '\\' && i+1 < len(raw) {
			i++
			continue
		}
		return false
	}
	return true
}
