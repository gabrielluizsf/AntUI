package css

// Gradient is one parsed gradient function. Angle is radians measured
// clockwise from the top, matching "to top" being 0°; Center anchors the
// radial or conic centre to a point in its painting area, with the default
// (built by the parsers) the middle. Shape and Size are radial-only.
type Gradient struct {
	Kind   uint8 // one of the Gradient* constants
	Angle  float64
	Center BackPos
	Shape  uint8
	Size   uint8
	Stops  []Stop
}
