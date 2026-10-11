package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// transformRegion is the part of the recorded layer worth compositing: the
// box widened by everything a box can shed beyond its edge — drops, spread,
// blur, outlines. The caller widens it further to what the widget actually
// drew, so content past the pad, like a select's menu, still comes over.
func (cs *cssStyle) transformRegion(st css.Style, x, y, w, h int) canvas.Area {
	pad := cs.transformPad(st)
	return canvas.Area{X: x - pad, Y: y - pad, Width: w + 2*pad, Height: h + 2*pad}
}
