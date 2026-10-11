package template

// tableWidth is how wide the table's content box is: the width the stylesheet
// gave it, or the width its columns ask for when it gave none — which is what a
// table does when nothing told it how wide to be. The width never exceeds the
// room the flow has for it.
func (c *CSS) tableWidth(b *tableBatch, cols []int, avail, chrome int) int {
	w := colsWidth(cols)
	if b.st.Has("width") {
		w = max(c.clampDim(b.st, c.style.length(b.st, b.st.Width, avail), "width", avail), 0)
	}
	return min(w, max(avail-chrome, 0))
}
