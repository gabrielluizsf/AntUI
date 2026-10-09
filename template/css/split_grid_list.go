package css

import "strings"

// splitGridList splits a track row into its space- or comma-separated tokens, keeping quoted strings and parenthesised groups whole.
func splitGridList(raw string, comma bool) []string {
	var out []string
	start, depth := -1, 0
	var quote byte
	flush := func(end int) {
		if start >= 0 {
			out = append(out, strings.TrimSpace(raw[start:end]))
			start = -1
			return
		}
		if comma {
			out = append(out, "")
		}
	}
	for i := 0; i < len(raw); i++ {
		ch := raw[i]
		if quote != 0 {
			if ch == '\\' {
				i++
			} else if ch == quote {
				quote = 0
			}
			continue
		}
		if ch == '"' || ch == '\'' {
			if start < 0 {
				start = i
			}
			quote = ch
			continue
		}
		switch ch {
		case '(':
			if start < 0 {
				start = i
			}
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',', ' ', '\t', '\n', '\r':
			if depth == 0 && (comma && ch == ',' || !comma && ch != ',') {
				flush(i)
			}
		default:
			if start < 0 {
				start = i
			}
		}
	}
	flush(len(raw))
	return out
}
