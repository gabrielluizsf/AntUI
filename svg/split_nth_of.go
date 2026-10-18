package svg

import "strings"

// splitNthOf separates an `:nth-…` argument into the `An+B` form and the `of S`
// list that may follow it. The `of` is a word on its own at the top level, so
// one written inside a string or a nested selector is not the one meant.
func splitNthOf(arg string) (form, of string) {
	depth := 0
	for i := 0; i < len(arg); i++ {
		switch arg[i] {
		case '(', '[':
			depth++
		case ')', ']':
			depth--
		case '"', '\'':
			if j := quoteEnd(arg, i); j >= 0 {
				i = j
			}
		case 'o', 'O':
			if depth == 0 && i > 0 && isSpace(arg[i-1]) && i+2 < len(arg) &&
				(arg[i+1] == 'f' || arg[i+1] == 'F') && isSpace(arg[i+2]) {
				return arg[:i], arg[i+2:]
			}
		}
	}
	return arg, ""
}

// stripNthSpace takes the whitespace out of an An+B piece, where it may stand
// around a sign, so "2n + 1" reads as "2n+1".
func stripNthSpace(s string) string {
	return strings.Join(strings.Fields(s), "")
}
