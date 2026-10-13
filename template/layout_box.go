package template

import "github.com/gabrielluizsf/antui/template/css"

// layoutBox is layout's engine: it resolves the entry's style for the given
// interaction state and places the widget it describes.
func (c *CSS) layoutBox(e *widgetEntry, role, label string, s State) (x, y, w, h int, st css.Style, ok bool) {
	st = c.style.entryStyle(e, role, s)
	if st.Display == css.DisplayNone {
		return 0, 0, 0, 0, st, false
	}
	u := c.style.u()
	winW, winH := c.win.Width(), c.win.Height()

	w, h = c.style.natural(role, label, st, u)
	if st.Has("width") {
		w = c.style.length(st, st.Width, winW)
	}
	if st.Has("max-width") {
		w = min(w, c.style.length(st, st.MaxWidth, winW))
	}
	if st.Has("min-width") {
		w = max(w, c.style.length(st, st.MinWidth, winW))
	}
	if st.Has("height") {
		h = c.style.length(st, st.Height, winH)
	}
	if st.Has("max-height") {
		h = min(h, c.style.length(st, st.MaxHeight, winH))
	}
	if st.Has("min-height") {
		h = max(h, c.style.length(st, st.MinHeight, winH))
	}

	ml := c.style.length(st, st.Margin[css.SideLeft], winW)
	mr := c.style.length(st, st.Margin[css.SideRight], winW)
	mt := c.style.length(st, st.Margin[css.SideTop], winW)
	mb := c.style.length(st, st.Margin[css.SideBottom], winW)

	if st.OutOfFlow() {
		x, y = c.placeOutOfFlow(st, w, h, winW, winH)
		return x, y, w, h, st, true
	}
	if st.Inline() {
		x, y = c.placeInline(w, h, winW, ml, mr)
	} else {
		x, y = c.placeBlock(w, h, winW, ml, mr, mt, mb)
	}
	x, y = c.shiftInset(st, x, y, winW, winH)
	return x, y, w, h, st, true
}
