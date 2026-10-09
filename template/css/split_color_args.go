package css

import (
	"strings"
)

// splitColorArgs cuts a functional colour's body ("1 2 3 / 0.5",
// "255, 0, 0, 0.5") into channel slots and, when a "/" introduced it, the
// trailing alpha. The comma form spreads its slots across the legacy syntax;
// the spaces form is the modern CSS4 way. ok is false when the trailing ")"
// is missing or the body is empty.
func splitColorArgs(s string) (slots []string, alpha string, hasAlpha, ok bool) {
	open := strings.IndexByte(s, '(')
	if open < 0 || !strings.HasSuffix(s, ")") {
		return nil, "", false, false
	}
	body := s[open+1 : len(s)-1]
	head, tail, slash := splitSlash(body)
	if strings.Contains(head, ",") {
		slots = splitFields(head, ',')
	} else {
		slots = strings.Fields(head)
	}
	for i := range slots {
		slots[i] = strings.TrimSpace(slots[i])
	}
	if slash {
		return slots, strings.TrimSpace(tail), true, true
	}
	return slots, "", false, true
}
