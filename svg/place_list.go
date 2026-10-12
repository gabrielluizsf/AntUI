package svg

// placeList puts one number on each letter of the writing. Every one of them
// has to be a number, since a list with a hole in it would put a letter in a
// place the file did not name, and the letters after it would follow it there.
func placeList(runs []TextRun, name, raw string, args []string, warn func(string, ...any)) {
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
