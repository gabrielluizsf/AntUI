package svg

import "strings"

// coord is one coordinate an element wrote as a single number, which is what
// `dx` and `dy` are: a move of the pen for the piece, never one for each
// letter. A number that is not one leaves the writing where it was, and a list
// of them is said rather than being read as though it were the first of them
// by accident.
func coord(e *element, name string, warn func(string, ...any)) (float64, bool) {
	raw := strings.TrimSpace(e.attr(name))
	if raw == "" {
		return 0, false
	}
	args := splitArgs(raw)
	if len(args) == 0 {
		return 0, false
	}
	v, err := parseNumber(args[0])
	if err != nil {
		warn("the %s %q is not a number, so the writing stays where it was", name, raw)
		return 0, false
	}
	if len(args) > 1 {
		warn("the %s gives a place for each letter of the writing, and only the first is used", name)
	}
	return v, true
}
