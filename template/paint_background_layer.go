package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintBackgroundLayer rasterizes one layer of the background — an image from
// the registry, or a gradient — into a fresh canvas the size of the box,
// clipped to the layer's own background-clip region. It returns nil when the
// layer has nothing to draw, and the canvas's clip bounds limit the paint in
// either case.
func (cs *cssStyle) paintBackgroundLayer(win *antui.Window, st css.Style, i, x, y, w, h int, cr float64) *canvas.Canvas {
	img := propAt(st.BackgroundImages, i, css.BackImage{})
	if img.Grad != nil {
		return cs.paintGradientLayer(st, i, *img.Grad, w, h, cr)
	}
	// A url() names either a canvas the developer registered on this template or
	// an icon of the icon set, and the registered canvas wins: it was named for
	// this template on purpose, where an icon is what a name falls back to. So a
	// developer who has registered something under a name that is also an icon
	// gets what they registered, and only the names nobody claimed reach the icon
	// set.
	if src := cs.image(img.URL); src != nil {
		return cs.paintImageLayer(win, st, i, src, w, h, cr)
	}
	if icon := cs.iconImage(img.URL, st); icon != nil {
		// Resolve background-size against the natural icon once, then ask the
		// SVG renderer for exactly that final size. paintImageLayer can then
		// composite it without a second bitmap resample.
		origin := propAt(st.BackgroundOrigin, i, css.BackPadding)
		_, _, ow, oh := cs.originRect(st, 0, 0, w, h, origin)
		dw, dh := cs.backgroundTileSize(st, i, icon, ow, oh)
		if direct := cs.iconImageAt(img.URL, st, dw, dh); direct != nil {
			icon = direct
		}
		return cs.paintImageLayer(win, st, i, icon, w, h, cr)
	}
	return nil
}
