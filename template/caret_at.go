package template

// caretAt maps a click at mx onto a caret position in text: each glyph's cell
// is split in half and snapped to its nearer edge, so a click lands where a
// text editor would put the caret. left is where the line's drawing starts
// and advance measures a text prefix exactly the way the box draws it.
// Positions before the first glyph return 0; past the last, len(text).
func caretAt(advance func(string) int, text string, left, mx int) int {
	if mx <= left {
		return 0
	}
	var (
		bounds []int // advance of each glyph boundary, in text order
		at     []int // byte index of each glyph boundary
	)
	for pos := 0; pos < len(text); {
		pos = utf8NextCursor(text, pos)
		bounds = append(bounds, advance(text[:pos]))
		at = append(at, pos)
	}
	for i, b := range bounds {
		prev := 0
		if i > 0 {
			prev = bounds[i-1]
		}
		edge := left + b
		if mx < edge {
			if i > 0 && mx < left+(prev+b)/2 {
				return at[i-1]
			}
			if i == 0 && mx < left+b/2 {
				return 0
			}
			return at[i]
		}
	}
	return len(text)
}
