package css

import "github.com/gabrielluizsf/antui/canvas"

// borderSides holds the border width, style and colour a shorthand spells
// out for the sides it names, plus which of the three it supplied at all.
type borderSides struct {
	w                [4]int
	s                [4]uint8
	c                [4]canvas.Color
	setW, setS, setC bool
}

// parseBorderSides reads the words of a border shorthand like "1px solid
// red", placing each part on the sides index names: sides[i] >= 0 selects
// that side, -1 marks a side the shorthand did not name.
func parseBorderSides(raw string, sides [4]int) borderSides {
	var b borderSides
	for _, wrd := range splitWords(raw) {
		if l, err := parseLength(wrd); err == nil && l.u == unitPx && !l.IsPct() {
			for i := 0; i < 4; i++ {
				if sides[i] >= 0 {
					b.w[sides[i]] = max(l.Px(0), 0)
				} else {
					b.w[i] = max(l.Px(0), 0)
				}
			}
			b.setW = true
			continue
		}
		if v, ok := parseBorderStyle(wrd); ok {
			for i := 0; i < 4; i++ {
				if sides[i] >= 0 {
					b.s[sides[i]] = v
				} else {
					b.s[i] = v
				}
			}
			b.setS = true
			continue
		}
		if col, err := ParseColor(wrd); err == nil {
			for i := 0; i < 4; i++ {
				if sides[i] >= 0 {
					b.c[sides[i]] = col
				} else {
					b.c[i] = col
				}
			}
			b.setC = true
		}
	}
	return b
}
