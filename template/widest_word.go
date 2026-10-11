package template

import (
	"strings"

	"github.com/gabrielluizsf/antui/template/css"
)

// widestWord is the width of the longest run of text the widget cannot break
// across two lines, measured with the style's own typography. Whitespace alone
// is no width at all.
func (cs *cssStyle) widestWord(st css.Style, label string) int {
	if label == "" {
		return 0
	}
	text := transformText(label, st.TextTransform)
	if !strings.ContainsAny(text, " \t\n\r") {
		return cs.measure(st, text)
	}
	widest := 0
	for _, word := range strings.Fields(text) {
		widest = max(widest, cs.measure(st, word))
	}
	return widest
}
