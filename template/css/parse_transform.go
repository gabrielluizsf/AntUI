package css

import (
	"strings"
)

// parseTransform reads a space-separated transform list: functions such as
// translate(10px 20%), rotate(45deg), scale(2), skew(20deg) or matrix(...).
// "none" clears the list, and any function the engine does not know drops the
// whole property exactly as an unknown value would.
func parseTransform(raw string, ctx Units) ([]TransformFunc, bool) {
	s := strings.TrimSpace(raw)
	if s == "none" {
		return nil, true
	}
	var out []TransformFunc
	for _, tk := range splitTokens(s) {
		open := strings.IndexByte(tk, '(')
		if open < 0 || !strings.HasSuffix(tk, ")") {
			return nil, false
		}
		name := strings.ToLower(strings.TrimSpace(tk[:open]))
		inner := strings.TrimSpace(tk[open+1 : len(tk)-1])
		f, ok := parseTransformFunc(name, inner, ctx)
		if !ok {
			return nil, false
		}
		out = append(out, f)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}
