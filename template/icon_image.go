package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// iconImage is the picture a url() names in the icon set, painted in the colour
// of the box's text, or nil when the name is not an icon.
//
// The size it is painted at is the icon's own size times the scale, so that
// `background-size: auto` — and anything that leaves the size alone — gives the
// icon at the size it was drawn to be on a screen of this scale, which is the
// same rule every other measurement in this package follows.
func (cs *cssStyle) iconImage(name string, st css.Style) *canvas.Canvas {
	if name == "" {
		return nil
	}
	iconMu.RLock()
	e := icons[name]
	iconMu.RUnlock()
	if e == nil || e.img == nil {
		return nil
	}
	w := int(e.img.Width * float64(cs.u()))
	h := int(e.img.Height * float64(cs.u()))
	if w <= 0 || h <= 0 {
		w, h = IconSize*cs.u(), IconSize*cs.u()
	}
	return IconCanvas(name, w, h, cs.ink(st))
}

// iconImageAt paints an icon at the size the background layer will use. This
// keeps SVG edges vector sharp: rendering at the viewBox's natural size and
// then resizing the bitmap would filter the same drawing a second time.
func (cs *cssStyle) iconImageAt(name string, st css.Style, w, h int) *canvas.Canvas {
	if name == "" || w <= 0 || h <= 0 {
		return nil
	}
	return IconCanvas(name, w, h, cs.ink(st))
}
