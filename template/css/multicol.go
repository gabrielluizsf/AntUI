package css

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// ColumnFill names how a multi-column container spends a height it was given.
// Balance, the initial value, spreads the content evenly so the columns come
// out the same height; auto fills the first column to the height before the
// next one starts.
const (
	ColumnFillBalance uint8 = iota
	ColumnFillAuto
)

// Break* names whether a box may be cut in the middle. The engine has one kind
// of fragmentation, the cut a multi-column container makes, and break-inside
// speaks to it: BreakAuto lets the flow cut a box wherever it runs out of
// room, BreakAvoid asks for the whole of it in one piece.
const (
	BreakAuto uint8 = iota
	BreakAvoid
)

// Float* names the side a box is taken out of the flow towards, and Clear* the
// sides a box asks to be pushed past. Both are read so the cascade has an
// answer for them; see doc.go for why this flow acts on neither.
const (
	FloatNone uint8 = iota
	FloatLeft
	FloatRight
	FloatInlineStart
	FloatInlineEnd
)

const (
	ClearNone uint8 = iota
	ClearLeft
	ClearRight
	ClearBoth
)

// parseColumnCount reads a column count. A whole number is the count; auto and
// an unreadable number leave it to the width, as the initial value does.
func parseColumnCount(raw string) (int, bool) {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "auto" {
		return 0, true
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, false
	}
	return n, true
}

// parseColumnWidth reads the width a column is given when the container is
// wide enough for it. A length is a width; auto leaves the count to decide.
func parseColumnWidth(raw string, ctx Units) (Length, bool) {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "auto" {
		return Auto(), true
	}
	l, err := parseLengthAt(s, ctx)
	if err != nil || l.Auto() || l.None() || l.value < 0 {
		return Auto(), false
	}
	return l, true
}

// parseColumns reads the columns shorthand, a count and a width in either
// order, which is how a stylesheet asks for "three columns, each at least this
// wide" in a single declaration. A lone number is a count, a lone length a
// width, and auto on its own changes neither.
func parseColumns(raw string, ctx Units) (count int, width Length, ok bool) {
	width = Auto()
	words := splitWords(raw)
	if len(words) == 0 {
		return 0, width, false
	}
	for _, wrd := range words {
		if strings.ToLower(wrd) == "auto" {
			continue
		}
		if n, err := strconv.Atoi(wrd); err == nil {
			if n < 0 {
				return 0, width, false
			}
			count = n
			continue
		}
		l, err := parseLengthAt(wrd, ctx)
		if err != nil || l.Auto() || l.None() || l.value < 0 {
			return 0, width, false
		}
		width = l
	}
	return count, width, true
}

// parseColumnFill reads how a declared height is spent across the columns.
func parseColumnFill(raw string) (uint8, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "balance":
		return ColumnFillBalance, true
	case "auto":
		return ColumnFillAuto, true
	}
	return 0, false
}

// parseBreakInside reads whether a box may be cut in the middle. The
// page-break-inside name the property had before break-inside replaced it means
// the same thing here, so both reach the same field. avoid-column, avoid-page
// and avoid-region all say avoid: the engine has one kind of fragment, so it
// has one answer to all of them.
func parseBreakInside(raw string) (uint8, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "auto":
		return BreakAuto, true
	case "avoid", "avoid-column", "avoid-page", "avoid-region":
		return BreakAvoid, true
	}
	return 0, false
}

// parseFloat reads the side a box is taken out of the flow towards.
func parseFloat(raw string) (uint8, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "none":
		return FloatNone, true
	case "left":
		return FloatLeft, true
	case "right":
		return FloatRight, true
	case "inline-start":
		return FloatInlineStart, true
	case "inline-end":
		return FloatInlineEnd, true
	}
	return 0, false
}

// parseClear reads the sides a box asks to be pushed past.
func parseClear(raw string) (uint8, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "none":
		return ClearNone, true
	case "left":
		return ClearLeft, true
	case "right":
		return ClearRight, true
	case "both":
		return ClearBoth, true
	}
	return 0, false
}

// applyColumnRule reads the column-rule shorthand, a width, a style and a
// colour in any order, and writes the three longhands. It takes the components
// it recognises the way the border shorthand does, and reports false for a
// value that names none of them so the caller can leave the rule alone.
func applyColumnRule(st *Style, set map[string]bool, raw string, note func()) bool {
	var (
		w                int
		s                uint8
		c                canvas.Color
		setW, setS, setC bool
	)
	words := splitWords(raw)
	if len(words) == 0 {
		return false
	}
	for _, wrd := range words {
		if l, err := parseLength(wrd); err == nil && l.u == unitPx {
			w = max(l.Px(0), 0)
			setW = true
			continue
		}
		if v, ok := parseBorderStyle(wrd); ok {
			s = v
			setS = true
			continue
		}
		if col, err := ParseColor(wrd); err == nil {
			c = col
			setC = true
		}
	}
	if setW {
		st.ColumnRuleWidth = w
		set["column-rule-width"] = true
	}
	if setS {
		st.ColumnRuleStyle = s
		set["column-rule-style"] = true
	}
	if setC {
		st.ColumnRuleColor = c
		set["column-rule-color"] = true
	}
	if setW || setS || setC {
		note()
	}
	return setW || setS || setC
}
