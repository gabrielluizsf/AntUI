package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// unPremultiply takes a colour out of the form the operations ran it in —
// held against the alpha that came out with it — and back to the straight
// colour a pixel is written as.
func unPremultiply(cr [3]float64, qr float64) canvas.Color {
	if qr <= 0 {
		return 0
	}
	return rgbaColor(math.Min(1, cr[0]/qr), math.Min(1, cr[1]/qr), math.Min(1, cr[2]/qr), qr)
}
