package template

import "github.com/gabrielluizsf/antui/template/css"

// closeFlex gives the batch its border-box height — away from a declared
// height, or from the content plus the frame — and records the content box
// and content origin that follow from it.
func (c *CSS) closeFlex(b *flexBatch, st css.Style) int {
	winH := c.win.Height()
	h := b.contentH + b.frameH
	if st.Has("height") {
		h = c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH)
	}
	h = max(h, c.clampDim(st, 0, "height", winH))
	b.h = h
	// The content box every child draws against is exactly what the border box
	// leaves inside, so an explicit height or a tall row agrees with the flow.
	b.contentH = max(h-b.frameH, 0)
	b.originY = b.y + b.frameTop
	return h
}
