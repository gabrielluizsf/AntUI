package svg

import (
	"fmt"
	"strings"
)

// parseRelativeList reads the inside of a `:has(…)`: one or more relative
// selectors, each optionally opening with the combinator that joins it to the
// element the `:has` stands on. With none written the lead is a descendant,
// which is what a relative selector with nothing in front means.
func parseRelativeList(s string, ns namespace) ([]relativeSelector, error) {
	var out []relativeSelector
	for _, part := range splitTopComma(s) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		lead := descendant
		if c := part[0]; c == '>' || c == '+' || c == '~' {
			lead = combinatorOf(c)
			part = strings.TrimSpace(part[1:])
		}
		cs, err := parseComplex(part, ns)
		if err != nil {
			return nil, err
		}
		out = append(out, relativeSelector{lead: lead, cs: cs})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("names no selector")
	}
	return out, nil
}
