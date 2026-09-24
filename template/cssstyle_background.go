package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// backgroundOn reports whether a box paints anything behind its border: a
// declared colour, a shorthand that touched the background, or any image
// layer. A box that answers false shows the window's own fill through it.
func (cs *cssStyle) backgroundOn(st css.Style) bool {
	if st.Has("background-color") || st.Has("background") {
		return true
	}
	return len(st.BackgroundImages) > 0
}

// paintBackground paints the surface of a box: the background colour first,
// then each layer above it — the last declared image sitting lowest, the
// first on top, the CSS stacking. A layer is rasterized into its own canvas
// and composited over the box so the alphas of several stacked layers aculate
// instead of overwriting one another, and each is cut at its background-clip
// box so a border-box layer does not bleed over the border.
func (cs *cssStyle) paintBackground(win *antui.Window, st css.Style, x, y, w, h, rx, ry int) {
	cv := win.Canvas()
	if st.Has("background-color") || st.Has("background") {
		cv.FillRoundRectXY(x, y, w, h, rx, ry, cs.fill(st))
	}
	n := len(st.BackgroundImages)
	if n == 0 {
		return
	}
	// Rounded corners are fractions of the box side, so a layer clipped to a
	// smaller box keeps the same curvature.
	cr := 0.0
	if min(w, h) > 0 {
		cr = max(float64(rx), float64(ry)) / float64(min(w, h))
	}
	for i := n - 1; i >= 0; i-- {
		layer := cs.paintBackgroundLayer(win, st, i, x, y, w, h, cr)
		if layer != nil {
			cv.BlitOver(x, y, layer)
		}
	}
}

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
	return cs.paintImageLayer(win, st, i, cs.image(img.URL), w, h, cr)
}

// propAt is the layer-list cycling: a single declared value applies to every
// layer, the next value to the next layer, wrapping around.
func propAt[T any](list []T, i int, def T) T {
	if len(list) == 0 {
		return def
	}
	return list[i%len(list)]
}

// originRect insets the box by the border for a padding box, or the border
// plus padding for a content box; the border box is the whole box. Each
// background attaches to its own positioning area.
func (cs *cssStyle) originRect(st css.Style, x, y, w, h int, origin uint8) (int, int, int, int) {
	u := cs.u()
	b := func(side int) int { return st.BorderWidth[side] * u }
	p := func(side int) int { return cs.base(st, st.Padding[side]) }
	switch origin {
	case css.BackBorder:
		return x, y, w, h
	case css.BackContent:
		return x + b(3) + p(3), y + b(0) + p(0),
			max(w-b(1)-b(3)-p(1)-p(3), 0), max(h-b(0)-b(2)-p(0)-p(2), 0)
	default: // the padding box is the initial background-origin
		return x + b(3), y + b(0), max(w-b(1)-b(3), 0), max(h-b(0)-b(2), 0)
	}
}

// maskLayer cuts an already-painted full-box layer down to the layer's
// background-clip region, keeping the box's corner curvature. A border-box
// clip masks the whole layer; a padding or content box moves the cut inside
// the border.
func (cs *cssStyle) maskLayer(layer *canvas.Canvas, st css.Style, clip uint8, x, y, w, h int, cr float64) {
	boxX, boxY, boxW, boxH := x, y, w, h
	switch clip {
	case css.BackPadding:
		boxX, boxY, boxW, boxH = cs.originRect(st, x, y, w, h, css.BackPadding)
	case css.BackContent:
		boxX, boxY, boxW, boxH = cs.originRect(st, x, y, w, h, css.BackContent)
	}
	if cr == 0 {
		layer.MaskRect(boxX, boxY, boxW, boxH)
		return
	}
	rr := min(boxW, boxH)
	layer.MaskRoundRect(boxX, boxY, boxW, boxH, int(cr*float64(rr)+0.5), int(cr*float64(rr)+0.5))
}

// autoDim reports whether a background-size length leaves the axis at the
// picture's own size: the auto keyword, the no-value zero Length, or any
// spelling of none.
func autoDim(l css.Length) bool {
	return l.Auto() || l.None() || l == (css.Length{})
}

// sizeDim measures one background-size axis against its positioning area,
// falling back onto the picture's own width or height.
func (cs *cssStyle) sizeDim(st css.Style, l css.Length, base int, natural int) int {
	if autoDim(l) {
		return natural
	}
	if l.IsPct() {
		return l.Resolve(css.Units{Width: base})
	}
	return l.Resolve(css.Units{Width: base}) * cs.u()
}

// edgePixel is where one edge of a picture lands along an axis of its
// positioning area: percentages measure the side left over, fixed offsets the
// area itself, and a right or bottom anchor folds its offset back towards the
// inside, as the CSS "right 10px" inset does.
func (cs *cssStyle) edgePixel(st css.Style, pos css.BackPosition, u, box, img int) int {
	left := 0
	base := max(box-img, 0)
	switch pos.Edge {
	case css.BackPosMiddle:
		left = base / 2
	case css.BackPosEnd:
		left = base
	}
	if pos.Off.IsPct() {
		off := pos.Off.Resolve(css.Units{Width: base})
		if pos.Edge == css.BackPosEnd {
			left -= off
		} else {
			left += off
		}
	} else {
		off := pos.Off.Resolve(css.Units{Width: box}) * u
		if pos.Edge == css.BackPosEnd {
			left -= off
		} else {
			left += off
		}
	}
	return left
}

// centerFrac is where the centre of a gradient lands as a fraction of one
// axis. A bare offset runs from the starting edge; a right or bottom anchor
// ranges back towards the inside.
func (cs *cssStyle) centerFrac(st css.Style, pos css.BackPosition, u, dim int) float64 {
	frac := 0.5
	switch pos.Edge {
	case css.BackPosEnd:
		frac = 1
	case css.BackPosStart, css.BackPosOffset:
		frac = 0
	}
	den := float64(max(dim, 1))
	var off float64
	if pos.Off.IsPct() {
		off = float64(pos.Off.Resolve(css.Units{Width: dim})) / den
	} else {
		off = float64(pos.Off.Resolve(css.Units{Width: dim})*u) / den
	}
	if pos.Edge == css.BackPosEnd {
		return frac - off
	}
	return frac + off
}

// coverSize is the smallest scale with the picture's proportions that still
// covers the box: the height scaled to the box height when that is the tight
// fit, the width otherwise.
func coverSize(srcW, srcH, boxW, boxH int) (int, int) {
	if srcW <= 0 || srcH <= 0 || boxW <= 0 || boxH <= 0 {
		return 0, 0
	}
	if srcW*boxH > boxW*srcH {
		return srcW * boxH / srcH, boxH
	}
	return boxW, srcH * boxW / srcW
}

// boxStart is the first tile position on an axis so the tiling lands on the
// requested edge, wrapping negative positions around the tile cycle.
func boxStart(p, dim int) int {
	if dim <= 0 {
		return 0
	}
	return p - ((p%dim)+dim)%dim
}

// paintGradientLayer rasterizes one gradient into a fresh box-sized canvas
// and clips it. The gradient family, shape and size constants of the engine
// coincide with the canvas's own, so the wire is a straight cast.
func (cs *cssStyle) paintGradientLayer(st css.Style, i int, g css.Gradient, w, h int, cr float64) *canvas.Canvas {
	u := cs.u()
	origin := propAt(st.BackgroundOrigin, i, css.BackPadding)
	ox, oy, ow, oh := cs.originRect(st, 0, 0, w, h, origin)
	if ow <= 0 || oh <= 0 {
		return nil
	}
	cg := canvas.Gradient{
		Kind:    canvas.GradientKind(uint8(g.Kind)),
		Angle:   g.Angle,
		CenterX: cs.centerFrac(st, g.Center[0], u, ow),
		CenterY: cs.centerFrac(st, g.Center[1], u, oh),
		Shape:   canvas.GradientShape(uint8(g.Shape)),
		Size:    canvas.GradientSize(uint8(g.Size)),
		Stops:   make([]canvas.GradientStop, 0, len(g.Stops)),
	}
	alpha := cs.alpha(st)
	for _, s := range g.Stops {
		a := int(s.Color.A()) * alpha / 255
		cg.Stops = append(cg.Stops, canvas.GradientStop{Offset: s.Offset, Color: canvas.Fade(s.Color, a)})
	}
	layer := boxLayer(w, h)
	if layer == nil {
		return nil
	}
	layer.FillGradient(ox, oy, ow, oh, cg)
	cs.maskLayer(layer, st, propAt(st.BackgroundClip, i, css.BackBorder), 0, 0, w, h, cr)
	return layer
}

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
	sz := propAt(st.BackgroundSize, i, css.BackSize{})
	dw, dh := src.Width, src.Height
	switch {
	case sz.Cover:
		dw, dh = coverSize(src.Width, src.Height, ow, oh)
	case sz.Contain:
		dw, dh = canvas.Fit(src.Width, src.Height, ow, oh)
	default:
		if !autoDim(sz.W) {
			dw = cs.sizeDim(st, sz.W, ow, src.Width)
		}
		if !autoDim(sz.H) {
			dh = cs.sizeDim(st, sz.H, oh, src.Height)
		}
	}
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
		layer.BlitScaled(px, py, dw, dh, src)
	} else {
		for tx := boxStart(px, dw); tx < w; tx += dw {
			for ty := boxStart(py, dh); ty < h; ty += dh {
				layer.BlitScaled(tx, ty, dw, dh, src)
			}
		}
	}
	cs.maskLayer(layer, st, propAt(st.BackgroundClip, i, css.BackBorder), 0, 0, w, h, cr)
	return layer
}

// boxLayer makes the empty canvas a background layer is rasterized into. A
// window-size canvas is the rare allocation that can fail; a nil result makes
// the caller skip the layer.
func boxLayer(w, h int) *canvas.Canvas {
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		return nil
	}
	return cv
}
