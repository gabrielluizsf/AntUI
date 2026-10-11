package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// drawInput paints a text field's box, text and caret.
func (win *Window) drawInput(s inputState, text *string) {
	const padding = 8

	win.cv.FillRoundRect(s.x, s.y, s.w, s.h, win.theme.Radius, win.theme.Surface)
	border := win.theme.Border
	switch {
	case s.focused:
		border = win.theme.Accent
	case s.hovered:
		border = win.theme.TextMuted
	}
	win.cv.RoundRect(s.x, s.y, s.w, s.h, win.theme.Radius, border)

	cursor := min(win.uiCursor, len(*text))
	cursorPx := canvas.TextWidth((*text)[:cursor])
	scroll := max(cursorPx-(s.w-2*padding), 0)

	savedClip := win.cv.Clip
	win.cv.SetClip(s.x+padding, s.y+1, s.w-2*padding, s.h-2)

	textY := s.y + (s.h-canvas.TextHeight())/2
	win.cv.Text(s.x+padding-scroll, textY, *text, win.theme.Text)

	if s.focused && int64(win.uiBlink*2)%2 == 0 {
		win.cv.FillRect(s.x+padding-scroll+cursorPx, textY-1, 1, canvas.TextHeight()+2, win.theme.Text)
	}

	win.cv.Clip = savedClip
}
