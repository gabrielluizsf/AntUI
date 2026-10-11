package canvas

// Dash is a stroke broken into pieces along its length, the dashes and the gaps
// between them measured in the same units as the path. A pattern of one number
// alternates that length on and the same length off; two numbers are a dash and
// a gap; an odd list is read twice over, which is how "ten on, five off" becomes
// dashes of ten with gaps of five and then gaps of ten with dashes of five. A
// list that is empty, or that holds a zero or a negative, is not a dash at all
// and the line comes back whole.
type Dash struct {
	On, Off []float64
	Offset  float64

	// The pattern as one run of dashes and gaps, and the distance a whole turn
	// of it measures, worked out once so that the walk does not build it again
	// for every point of every subpath.
	pattern []float64
	period  float64
}
