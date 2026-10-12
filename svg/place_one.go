package svg

// placeOne puts one number on the whole of the writing.
func placeOne(runs []TextRun, name, raw, arg string, warn func(string, ...any)) {
	v, err := parseNumber(arg)
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
}
