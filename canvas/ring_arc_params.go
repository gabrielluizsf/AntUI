package canvas

// RingArcParams holds parameters for drawing a thick ring arc.
type RingArcParams struct {
	Center    Point
	Radius    int
	Thickness int
	Start     float64
	Sweep     float64
	Color     Color
}
