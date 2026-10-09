package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// FontFace is one @font-face block: the family, weight and slant it was
// declared for, and the file it read. A face is read while the sheet parses,
// so a frame only ever looks one up.
type FontFace struct {
	Family string // lowercased, the way font-family names match
	Weight uint16 // the weight it declared, on the hundreds CSS matches on
	Style  uint8  // one of the FontStyle* constants
	face   *canvas.Face
}
