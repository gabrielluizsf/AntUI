package canvas

// faceWidth sums one glyph at a time so the per-glyph and per-space spacing
// land exactly where Draw puts them.
func faceWidth(f *Face, text string, o TextStyle) int {
	o = o.norm()
	fs := f.Scaled(o.Scale)
	widest, line := 0, 0
	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, line)
			line = 0
		case '\r':
		case '\t':
			line += fs.Width("\t")
		case ' ':
			line += fs.Width(" ") + o.WordSpacing
		default:
			line += fs.Width(string(r)) + o.LetterSpacing
		}
	}
	return max(widest, line)
}
