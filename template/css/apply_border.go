package css

// applyBorder sets one or all sides of the border from a shorthand like
// "1px solid red". sides[i] >= 0 names the side to set, or all four.
func applyBorder(st *Style, set map[string]bool, sides [4]int, raw string, note func()) {
	b := parseBorderSides(raw, sides)
	for i := 0; i < 4; i++ {
		dst := i
		if sides[i] >= 0 {
			dst = sides[i]
		}
		if b.setW {
			st.BorderWidth[dst] = b.w[dst]
		}
		if b.setS {
			st.BoxStyle[dst] = b.s[dst]
		}
		if b.setC {
			st.BoxColor[dst] = b.c[dst]
		}
	}
	if b.setW {
		set["border-width"] = true
	}
	if b.setS {
		set["border-style"] = true
	}
	if b.setC {
		set["border-color"] = true
	}
	note()
}
