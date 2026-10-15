package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// drawInput paints a text field's box, text and caret.
func (win *Window) drawInput(s inputState, text *string) {
	const padding = 8

	win.cv.FillRoundRect(canvas.Round(canvas.Area{X: s.x, Y: s.y, Width: s.w, Height: s.h}, win.theme.Radius), win.theme.Surface)
	border := win.theme.Border
	switch {
	case s.focused:
		border = win.theme.Accent
	case s.hovered:
		border = win.theme.TextMuted
	}
	win.cv.RoundRect(canvas.Round(canvas.Area{X: s.x, Y: s.y, Width: s.w, Height: s.h}, win.theme.Radius), border)

	cursor := min(win.uiCursor, len(*text))
	cursorPx := canvas.TextWidth((*text)[:cursor])
	scroll := max(cursorPx-(s.w-2*padding), 0)

	savedClip := win.cv.Clip
	win.cv.SetClip(canvas.Area{X: s.x + padding, Y: s.y + 1, Width: s.w - 2*padding, Height: s.h - 2})

	textY := s.y + (s.h-canvas.TextHeight())/2
	win.cv.Text(canvas.Point{X: float64(s.x + padding - scroll), Y: float64(textY)}, *text, win.theme.Text)

	if s.focused && int64(win.uiBlink*2)%2 == 0 {
		win.cv.FillRect(canvas.Area{X: s.x + padding - scroll + cursorPx, Y: textY - 1, Width: 1, Height: canvas.TextHeight() + 2}, win.theme.Text)
	}

	win.cv.Clip = savedClip
}
