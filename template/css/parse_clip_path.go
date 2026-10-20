package css

import "strings"

// parseClipPath reads a clip-path value. Only the polygon form is used by the
// templates that reshape a widget; "none" clears the clip, and every other
// shape the engine cannot draw leaves the box unclipped rather than failing
// the declaration. Vertices come back in declaration order, each a pair of
// box-relative coordinate expressions the painter resolves against the box.
func parseClipPath(raw string) ([][2]string, bool) {
	s := strings.TrimSpace(strings.ToLower(raw))
	if s == "none" {
		return nil, true
	}
	open := strings.IndexByte(s, '(')
	if open < 0 || !strings.HasSuffix(s, ")") {
		return nil, false
	}
	if keyword(s[:open]) != "polygon" {
		return nil, false
	}
	inner := strings.TrimSpace(s[open+1 : len(s)-1])
	// A fill rule ("nonzero" / "evenodd") may lead the vertex list.
	if sp := strings.IndexAny(inner, " ,"); sp > 0 {
		if head := keyword(inner[:sp]); head == "nonzero" || head == "evenodd" {
			inner = strings.TrimSpace(inner[sp:])
		}
	}
	var out [][2]string
	for _, pair := range splitTopLevel(inner) {
		x, y, ok := splitPair(pair)
		if !ok {
			return nil, false
		}
		out = append(out, [2]string{x, y})
	}
	if len(out) < 3 {
		return nil, false
	}
	return out, true
}

// splitPair cuts a clip-path vertex into its x and y expressions at the first
// top-level space, so a space inside calc() is not mistaken for the separator.
func splitPair(vertex string) (string, string, bool) {
	depth := 0
	for i := 0; i < len(vertex); i++ {
		switch vertex[i] {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ' ', '\t':
			if depth == 0 {
				x := strings.TrimSpace(vertex[:i])
				y := strings.TrimSpace(vertex[i:])
				if x == "" || y == "" {
					return "", "", false
				}
				return x, y, true
			}
		}
	}
	return "", "", false
}
