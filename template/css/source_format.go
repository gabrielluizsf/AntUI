package css

import (
	"strings"
)

// sourceFormat reads the format("…") clause of a src source, lowercased, or
// "" when the source carries none. It is a hint only: a format this engine
// reads is believed, and one it does not is the only reason to pass a source
// over without opening it.
func sourceFormat(part string) string {
	lower := strings.ToLower(part)
	i := strings.Index(lower, "format(")
	if i < 0 {
		return ""
	}
	s := strings.TrimSpace(lower[i+len("format("):])
	end := strings.IndexByte(s, ')')
	if end < 0 {
		return ""
	}
	return unquote(strings.TrimSpace(s[:end]))
}
