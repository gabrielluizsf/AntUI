package svg

// firstStrongDir reads the direction an `auto` element is given: the way the
// first strong directional character in the writing it holds runs. The writing
// beside an element counts before the writing inside it, the same order the
// file wrote them in.
func firstStrongDir(e *element) string {
	if d := firstStrongIn(e.Text); d != "" {
		return d
	}
	for _, p := range e.Parts {
		if p.Elem == nil {
			if d := firstStrongIn(p.Text); d != "" {
				return d
			}
			continue
		}
		if d := firstStrongDir(p.Elem); d != "" {
			return d
		}
	}
	return "ltr"
}

// firstStrongIn reads the direction of the first character in the writing that
// runs one way or the other on its own: a letter of the Latin, Greek, Cyrillic
// or the like runs left to right, one of the Hebrew, Arabic or the like runs
// right to left, and everything between them — digits, spaces, punctuation —
// says nothing and is passed over.
func firstStrongIn(s string) string {
	for _, r := range s {
		if isStrongRTL(r) {
			return "rtl"
		}
		if isStrongLTR(r) {
			return "ltr"
		}
	}
	return ""
}
