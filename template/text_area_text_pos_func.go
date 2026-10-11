package template

// textAreaTextPos maps a click inside a wrapped multi-line text box onto the
// byte index of the character it lands on. It replicates where the box's
// drawing puts each line, so the caret lands exactly under the pointer.
// padding and lineH describe the box's inside margin and row height in
// window pixels, advance measures a prefix the way the box draws, and mx,my
// are the click in window coordinates.
func textAreaTextPos(text string, lines [][2]int, padding, x, y, h, cursor int, lineH int, advance func(string) int, mx, my int) int {
	visible := max((h-2*padding)/lineH, 1)
	viewTop := textViewTop(lines, cursor, visible)
	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := y + padding + (lif-viewTop)*lineH
		if my >= cy && my < cy+lineH {
			return l[0] + caretAt(advance, text[l[0]:l[1]], x+padding, mx)
		}
	}
	return len(text)
}
