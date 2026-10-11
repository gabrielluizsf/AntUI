package template

import "github.com/gabrielluizsf/antui/canvas"

// fitText truncates a run with an ellipsis so it fits maxW, working rune by
// rune so the cut never splits a character.
func fitText(cv *canvas.Canvas, text string, o canvas.TextStyle, maxW int) string {
	if maxW <= 0 || cv.TextWidthStyled(text, o) <= maxW {
		return text
	}
	const dot = "…"
	dw := cv.TextWidthStyled(dot, o)
	if dw > maxW {
		return ""
	}
	runes := []rune(text)
	for i := len(runes) - 1; i > 0; i-- {
		if cv.TextWidthStyled(string(runes[:i])+dot, o) <= maxW {
			return string(runes[:i]) + dot
		}
	}
	return dot
}
