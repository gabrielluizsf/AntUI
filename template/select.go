package template

import "github.com/gabrielluizsf/antui"

func (builtinStyle) Select(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	local := th.Surface
	if s.Hovered {
		local = th.SurfaceHover
	}
	border := th.Border
	if s.Focused || open {
		border = th.Accent
	}
	cv.FillRoundRect(x, y, w, h, th.Radius, local)
	cv.RoundRect(x, y, w, h, th.Radius, border)

	padding := 8 * u
	textY := y + (h-textHeight(u))/2
	cv.TextScaled(x+padding, textY, value, th.Text, u)

	// The arrow that says "there is more below". It is drawn inside a square
	// bigger than the arrow so that the arrow itself comes out at the size the
	// drawing says it is, centred in the square whatever the scale.
	caretBox := 10 * u
	DrawIcon(cv, IconCaret,
		x+w-padding-3*u-caretBox/2, y+(h-caretBox)/2, caretBox, caretBox, th.TextMuted)
}
