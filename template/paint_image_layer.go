package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintImageLayer rasterizes one url() layer: the picture is resized to its
// background-size, anchored by its background-position inside the
// background-origin area, and repeated when asked to. Each tile is scaled into
// the layer canvas, so a source with straight alpha over paints the layer the
// way CSS lays an image over the colour beneath it.
func (cs *cssStyle) paintImageLayer(win *antui.Window, st css.Style, i int, src *canvas.Canvas, w, h int, cr float64) *canvas.Canvas {
	if src == nil || src.Width <= 0 || src.Height <= 0 {
		return nil
	}
	u := cs.u()
	origin := propAt(st.BackgroundOrigin, i, css.BackPadding)
	ox, oy, ow, oh := cs.originRect(st, 0, 0, w, h, origin)
	if ow <= 0 || oh <= 0 {
		return nil
	}
	dw, dh := cs.backgroundTileSize(st, i, src, ow, oh)
	if dw <= 0 || dh <= 0 {
		return nil
	}
	pos := propAt(st.BackgroundPos, i, css.BackPos{})
	px := ox + cs.edgePixel(st, pos[0], u, ow, dw)
	py := oy + cs.edgePixel(st, pos[1], u, oh, dh)
	layer := boxLayer(w, h)
	if layer == nil {
		return nil
	}
	repeat := propAt(st.BackgroundRepeat, i, css.BackRepeat{css.BackRepeatRepeat, css.BackRepeatRepeat})
	if repeat[0] == css.BackRepeatNoRepeat && repeat[1] == css.BackRepeatNoRepeat {
		layer.BlitScaled(canvas.Area{X: px, Y: py, Width: dw, Height: dh}, src)
	} else {
		blitTiled(layer, src, px, py, dw, dh, w, h)
	}
	cs.maskLayer(layer, st, propAt(st.BackgroundClip, i, css.BackBorder), 0, 0, w, h, cr)
	return layer
}
