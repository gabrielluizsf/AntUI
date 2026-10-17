package svg

import "strings"

// parseSelectorList reads the selector side of a rule — one or more selectors
// separated by commas — into the selectors it names.
func parseSelectorList(s string) ([]complexSelector, error) {
	var out []complexSelector
	for _, part := range splitTopComma(s) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		cs, err := parseComplex(part)
		if err != nil {
			return nil, err
		}
		out = append(out, cs)
	}
	return out, nil
}

// splitTopComma splits on the commas that separate selectors, and not on a
// comma inside an attribute selector's brackets.
func splitTopComma(s string) []string {
	var out []string
	depth, start := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '[':
			depth++
		case ']':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	return append(out, s[start:])
}
