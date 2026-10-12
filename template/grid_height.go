package template

import "github.com/gabrielluizsf/antui/template/css"

// gridHeight applies the container's min-height and max-height to the content
// it solved, then its declared height, and clamps the result.
func (c *CSS) gridHeight(st css.Style, b *gridBatch, top, bottom int) int {
	winH := c.win.Height()
	if st.Has("min-height") && gridLengthDefinite(st.MinHeight) {
		limit := c.clampDim(st, c.style.length(st, st.MinHeight, winH), "height", winH)
		b.contentH = max(b.contentH, max(limit-top-bottom, 0))
	}
	if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) {
		limit := c.clampDim(st, c.style.length(st, st.MaxHeight, winH), "height", winH)
		b.contentH = min(b.contentH, max(limit-top-bottom, 0))
	}
	h := b.contentH + top + bottom
	if st.Has("height") && gridLengthDefinite(st.Height) {
		h = c.style.length(st, st.Height, winH)
	}
	return c.clampDim(st, h, "height", winH)
}
