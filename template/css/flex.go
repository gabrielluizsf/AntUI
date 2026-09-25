package css

import (
	"strconv"
)

// FlexDirection names the main axis of a flex container: which way its items
// pack. Row lays them left to right, column top to bottom, and the reverse
// variants pack from the far end of that axis.
const (
	FlexDirectionRow uint8 = iota
	FlexDirectionRowReverse
	FlexDirectionColumn
	FlexDirectionColumnReverse
)

// FlexWrap names how flex items give the main axis if they overflow it:
// nowrap lets them shrink and stick to one line, wrap grants lines, and
// wrap-reverse stacks the wrapped lines from the far end of the cross axis.
const (
	FlexWrapNowrap uint8 = iota
	FlexWrapWrap
	FlexWrapWrapReverse
)

// JustifyContent names how leftover main-axis space is spread between items:
// packed at one end, centred, or turned into spacing between, around or
// even around the items. flex-start is also the initial value.
const (
	JustifyFlexStart uint8 = iota
	JustifyFlexEnd
	JustifyCenter
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// Align* names how an item is placed across the cross axis of its line.
// AlignAuto is only legal for align-self, where it means "use the container's
// align-items"; AlignStretch is align-items' initial value, which grows an
// auto-sized item to fill its line's cross size. AlignBaseline is accepted
// and drawn as flex-start: the engine has no shared baseline model.
const (
	AlignAuto    uint8 = 0xFF
	AlignStretch uint8 = iota
	AlignFlexStart
	AlignFlexEnd
	AlignCenter
	AlignBaseline
)

// Content* name how leftover cross-axis space is distributed between whole
// lines when a container wraps. Stretch (the initial value) grows the lines,
// the rest spread or pack them.
const (
	ContentStretch uint8 = iota
	ContentFlexStart
	ContentFlexEnd
	ContentCenter
	ContentSpaceBetween
	ContentSpaceAround
	ContentSpaceEvenly
)

// parseFlexDirection turns "row", "column" and their reverses into a
// FlexDirection* constant.
func parseFlexDirection(raw string) (uint8, bool) {
	switch raw {
	case "row":
		return FlexDirectionRow, true
	case "row-reverse":
		return FlexDirectionRowReverse, true
	case "column":
		return FlexDirectionColumn, true
	case "column-reverse":
		return FlexDirectionColumnReverse, true
	}
	return 0, false
}

// parseFlexWrap turns "nowrap", "wrap" and "wrap-reverse" into a FlexWrap*
// constant.
func parseFlexWrap(raw string) (uint8, bool) {
	switch raw {
	case "nowrap":
		return FlexWrapNowrap, true
	case "wrap":
		return FlexWrapWrap, true
	case "wrap-reverse":
		return FlexWrapWrapReverse, true
	}
	return 0, false
}

// parseJustifyContent turns the six justify-content keywords into a
// Justify* constant. The bare start/end keywords behave like their
// flex- counterparts.
func parseJustifyContent(raw string) (uint8, bool) {
	switch raw {
	case "flex-start", "start":
		return JustifyFlexStart, true
	case "flex-end", "end":
		return JustifyFlexEnd, true
	case "center":
		return JustifyCenter, true
	case "space-between":
		return JustifySpaceBetween, true
	case "space-around":
		return JustifySpaceAround, true
	case "space-evenly":
		return JustifySpaceEvenly, true
	}
	return 0, false
}

// parseAlignItems turns the align-items keywords into an Align* constant;
// normal is stretch, the initial value.
func parseAlignItems(raw string) (uint8, bool) {
	switch raw {
	case "normal", "stretch":
		return AlignStretch, true
	case "flex-start", "start", "self-start":
		return AlignFlexStart, true
	case "flex-end", "end", "self-end":
		return AlignFlexEnd, true
	case "center":
		return AlignCenter, true
	case "baseline":
		return AlignBaseline, true
	}
	return 0, false
}

// parseAlignSelf turns the align-self keywords into an Align* constant,
// auto included — the only keyword align-items rejects.
func parseAlignSelf(raw string) (uint8, bool) {
	if raw == "auto" {
		return AlignAuto, true
	}
	return parseAlignItems(raw)
}

// parseAlignContent turns the align-content keywords into a Content*
// constant; normal behaves like stretch.
func parseAlignContent(raw string) (uint8, bool) {
	switch raw {
	case "normal", "stretch":
		return ContentStretch, true
	case "flex-start", "start":
		return ContentFlexStart, true
	case "flex-end", "end":
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

// parseGap reads a gap shorthand: one length for both axes, or a row gap
// followed by a column gap. A negative length or an unreadable token drops
// the whole declaration.
func parseGap(raw string, ctx Units) (row, col Length, ok bool) {
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 2 {
		return row, col, false
	}
	vals := make([]Length, len(parts))
	for i, p := range parts {
		if p == "normal" {
			vals[i] = Zero()
			continue
		}
		l, err := parseLengthAt(p, ctx)
		if err != nil || l.Auto() || l.None() || l.value < 0 {
			return row, col, false
		}
		vals[i] = l
	}
	if len(vals) == 1 {
		return vals[0], vals[0], true
	}
	return vals[0], vals[1], true
}

// parseOrder reads an integer order, auto resets it to zero.
func parseOrder(raw string) (int, bool) {
	if raw == "auto" {
		return 0, true
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return v, true
}

// parseFlexNumber reads a non-negative flex factor such as the grow or
// shrink multiplier; a negative number is invalid, as in CSS.
func parseFlexNumber(raw string) (float64, bool) {
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || v < 0 {
		return 0, false
	}
	return v, true
}

// parseFlexBasis reads the flex-basis longhand: a length, a percentage,
// auto, or content. content maps to the none marker so the solver reads it
// as "measure the item's own main size", the same way auto is handled.
func parseFlexBasis(raw string, ctx Units) (Length, bool) {
	if raw == "auto" {
		return Auto(), true
	}
	if raw == "content" || raw == "max-content" || raw == "min-content" || raw == "fit-content" {
		return Length{u: unitNone}, true
	}
	l, err := parseLengthAt(raw, ctx)
	if err != nil {
		return l, false
	}
	return l, true
}

// parseFlex reads the flex shorthand: none, auto, initial, or up to three
// values holding grow, shrink and basis in any order. A unitless number is
// the next-flex factor (grow first, shrink second); a length or a keyword is
// the basis. What the declaration does not say defaults to grow 1, shrink 1,
// basis 0, so "flex: 2" and "flex: 2 1 0%" agree.
func parseFlex(raw string, ctx Units) (grow, shrink float64, basis Length, ok bool) {
	switch raw {
	case "none":
		return 0, 0, Auto(), true
	case "auto":
		return 1, 1, Auto(), true
	case "initial":
		return 0, 1, Auto(), true
	}
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 3 {
		return 0, 0, Length{}, false
	}
	seenGrow, seenShrink := false, false
	hasBasis := false
	for _, p := range parts {
		if v, num := parseFlexNumber(p); num && !seenShrink {
			if !seenGrow {
				grow, seenGrow = v, true
				continue
			}
			shrink, seenShrink = v, true
			continue
		}
		if b, good := parseFlexBasis(p, ctx); good && !hasBasis {
			basis, hasBasis = b, true
			continue
		}
		return 0, 0, Length{}, false
	}
	if !seenGrow {
		grow = 1
	}
	if !seenShrink {
		shrink = 1
	}
	if !hasBasis {
		basis = Pct(0)
	}
	return grow, shrink, basis, true
}
