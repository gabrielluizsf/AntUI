package template

import "github.com/gabrielluizsf/antui"

func (builtinStyle) DatePickerBox(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
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
	cv.TextScaled(x+padding, y+(h-textHeight(u))/2, value, th.Text, u)

	// The caret that says "there is a calendar below", flipped the way the
	// dropdown's is while the picker is open.
	name := IconCaret
	if open {
		name = IconCaretUp
	}
	caretBox := 10 * u
	DrawIcon(cv, name,
		x+w-padding-3*u-caretBox/2, y+(h-caretBox)/2, caretBox, caretBox, th.TextMuted)
}
