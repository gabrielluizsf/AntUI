package svg

import "strings"

// placeAt reads one of the two names a place may be written under, either as
// one number for the whole of the writing or as one number for each letter of
// it.
func placeAt(runs []TextRun, e *element, name string, warn func(string, ...any)) {
	raw := strings.TrimSpace(e.attr(name))
	if raw == "" {
		return
	}
	args := splitArgs(raw)
	if len(args) == 0 {
		return
	}
	if len(args) == 1 {
		v, err := parseNumber(args[0])
		if err != nil {
			warn("the %s %q is not a number, so the writing stays where it was", name, raw)
			return
		}
		at := &runs[0]
		if name == "y" {
			if !at.HasY {
				at.Y, at.HasY = v, true
			}
			return
		}
		if !at.HasX {
			at.X, at.HasX = v, true
		}
		return
	}
	// One number for each letter: every one of them has to be a number, since
	// a list with a hole in it would put a letter in a place the file did not
	// name, and the letters after it would follow it there.
	vals := make([]float64, len(args))
	for i, a := range args {
		v, err := parseNumber(a)
		if err != nil {
			warn("the %s %q is not a list of numbers, so the writing stays where it was", name, raw)
			return
		}
		vals[i] = v
	}
	seen := 0
	for i := range runs {
		if seen >= len(vals) {
			return
		}
		letters := len([]rune(runs[i].Text))
		if letters == 0 {
			continue
		}
		take := min(len(vals)-seen, letters)
		placed := runs[i].HasX
		if name == "y" {
			placed = runs[i].HasY
		}
		if !placed {
			if name == "y" {
				runs[i].YList = vals[seen : seen+take]
			} else {
				runs[i].XList = vals[seen : seen+take]
			}
		}
		seen += letters
	}
}
