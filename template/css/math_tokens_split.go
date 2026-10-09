package css

import (
	"strings"
)

func (t *mathTokens) split() []string {
	if t.toks != nil {
		return t.toks
	}
	s := t.s
	var out []string
	for i := 0; i < len(s); {
		c := s[i]
		switch {
		case isSpace(c):
			i++
		case strings.IndexByte("+-*/(),", c) >= 0:
			// A function name ("calc", "min"… ) is the letters just before
			// its "(", handed to the parser as one token.
			if c == '(' {
				name := ""
				j := i - 1
				for j >= 0 && isIdentByte(s[j]) {
					name = string(s[j]) + name
					j--
				}
				switch name {
				case "calc", "min", "max", "clamp":
					out = append(out, name)
					out = append(out, "(")
					i++
					continue
				}
			}
			out = append(out, string(c))
			i++
		case c == '.' || (c >= '0' && c <= '9'):
			j := i
			for j < len(s) && (s[j] == '.' || (s[j] >= '0' && s[j] <= '9')) {
				j++
			}
			for j < len(s) && isUnitByte(s[j]) {
				j++
			}
			out = append(out, s[i:j])
			i = j
		default:
			// A stray glyph: skip it rather than fail the whole formula.
			i++
		}
	}
	t.toks = out
	return out
}
