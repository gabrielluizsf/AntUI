package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// strokeWriting draws what a stroke over the writing covers: the same writing
// at every point around a circle as wide across as the stroke, which is the
// shape a stroke of that width makes around letters. The circle is walked
// finely enough that its points land about half a pixel apart, and a stroke
// narrower than a pixel is taken at half a pixel rather than at nothing, so a
// hairline still shows after the picture is brought back down to size.
//
// The outline comes out round rather than at the miter a stroke takes round the
// corner of a shape: the corners a glyph has are smaller than the difference
// would be at the sizes writing is drawn at.
func strokeWriting(cv *canvas.Canvas, put putRun, width float64, c canvas.Color, scale float64) {
	r := math.Max(width/2, 0.5)
	n := min(max(int(math.Ceil(4*math.Pi*r)), 8), 128)
	for i := range n {
		a := 2 * math.Pi * float64(i) / float64(n)
		put(cv, c, math.Cos(a)*r/scale, math.Sin(a)*r/scale)
	}
}
