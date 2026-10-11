package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// LabelHeight is how tall one is with pad above and below.
func LabelHeight(pad int) int { return canvas.FontHeight + pad*2 }
