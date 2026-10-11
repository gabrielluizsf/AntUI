package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// clipRuleOf is one value of it.
func clipRuleOf(raw string) canvas.FillRule {
	if strings.EqualFold(strings.TrimSpace(raw), "evenodd") {
		return canvas.FillEvenOdd
	}
	return canvas.FillNonZero
}
