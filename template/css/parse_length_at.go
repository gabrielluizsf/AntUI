package css

import (
	"fmt"
	"strings"
)

// parseLengthAt reads a length with a measurement context, so % and the
// viewport/font units resolve and calc() min() max() clamp() can evaluate.
// A calc() that balances is converted out of its formula into plain px.
func parseLengthAt(raw string, ctx Units) (Length, error) {
	s := strings.TrimSpace(strings.ToLower(raw))
	if s == "" {
		return Length{}, fmt.Errorf("css: empty length")
	}
	switch s {
	case "auto":
		return Auto(), nil
	case "none":
		return Length{u: unitNone}, nil
	case "inherit", "initial", "unset":
		return Length{u: unitNone}, nil
	}
	if mathLook(s) {
		if px, ok := EvalMath(s, ctx); ok {
			return Fixed(px), nil
		}
		return Length{}, fmt.Errorf("css: %q does not evaluate", raw)
	}
	return parseLength(s)
}
