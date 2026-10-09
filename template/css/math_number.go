package css

import (
	"strings"
)

import (
	"strconv"
)

import (
	"fmt"
)

// number reads a numeric literal with an optional unit, converted to
// reference pixels under the parser's context. A plain number is pixels.
func (p *mathParser) number() (float64, error) {
	tok := p.t.next()
	if tok == "" {
		return 0, fmt.Errorf("css: expected a number")
	}
	// The longest matching unit wins — rem before em, vmin before vw, etc.
	u := unitPx
	for _, suf := range unitSuffixes {
		if strings.HasSuffix(tok, suf.name) {
			u = suf.u
			tok = strings.TrimSuffix(tok, suf.name)
			break
		}
	}
	if strings.HasSuffix(tok, "%") {
		u = unitPct
		tok = strings.TrimSuffix(tok, "%")
	}
	v, err := strconv.ParseFloat(tok, 64)
	if err != nil {
		return 0, fmt.Errorf("css: %q is not a number", tok)
	}
	l := Length{u: u, value: v}
	return l.ref(p.ctx), nil
}
