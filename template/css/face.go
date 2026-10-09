package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Face is the font file itself, at the size the program draws its text at,
// with the names after it in the font-family list behind it for the runes it
// has no shape for.
func (ff *FontFace) Face() *canvas.Face { return ff.face }
