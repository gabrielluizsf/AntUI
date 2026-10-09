package css

func setLenSide(sides []Length, canonical string, set map[string]bool, i int, raw string, ctx Units) {
	l, err := parseLengthAt(raw, ctx)
	if err != nil {
		return
	}
	sides[i] = l
	set[canonical] = true
}
