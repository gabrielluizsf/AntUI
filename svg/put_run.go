package svg

import "github.com/gabrielluizsf/antui/canvas"

// putRun is one piece of the writing laid down in a canvas: in the colour it
// is given, moved by however far the caller asks. The fill asks for no move at
// all, and the stroke asks for one for each point around the circle its width
// makes — which is the whole of what a stroke of that width covers.
type putRun func(target *canvas.Canvas, c canvas.Color, dx, dy float64)
