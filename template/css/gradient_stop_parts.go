package css

import "strings"

// gradientStopParts returns the colour-stop arguments after a gradient
// prelude. Words that never became prelude are re-joined onto the front, so a
// colour stop that was missing its comma still reads as one argument.
func gradientStopParts(tokens []string, i int, tail []string) []string {
	if i < len(tokens) {
		return append([]string{strings.Join(tokens[i:], " ")}, tail...)
	}
	return tail
}
