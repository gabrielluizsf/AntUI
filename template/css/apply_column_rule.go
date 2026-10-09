package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

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
