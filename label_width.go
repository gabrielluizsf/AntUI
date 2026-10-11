package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// LabelWidth is how wide a piece of text is with pad either side of it.
func LabelWidth(text string, pad int) int {
	return canvas.TextWidth(text) + pad*2
}
