package css

import (
	"strings"
)

func splitTopLevel(text string) []string {
	var out []string
	depth := 0
	start := 0
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '(', '[':
			depth++
		case ')', ']':
			if depth > 0 {
				depth--
			}
		case '"', '\'':
			if end := quoteEnd(text, i); end >= 0 {
				i = end
			}
		case ',':
			if depth == 0 {
				out = append(out, text[start:i])
				start = i + 1
			}
		}
	}
	out = append(out, text[start:])
	trimmed := out[:0]
	for _, o := range out {
		if t := strings.TrimSpace(o); t != "" {
			trimmed = append(trimmed, t)
		}
	}
	return trimmed
}
