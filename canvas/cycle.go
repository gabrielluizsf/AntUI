package canvas

// pattern is the dash pattern as one run, dashes and gaps in turn, and period
// how far it has to go before it starts again. Both are worked out once, because
// every step of the walk asks for them.
func (d Dash) cycle() ([]float64, float64) {
	if len(d.pattern) > 0 {
		return d.pattern, d.period
	}
	// A pattern built by hand rather than by NewDash still has to be read.
	out := make([]float64, 0, len(d.On)*2)
	for i, l := range d.On {
		out = append(out, l, d.Off[i%len(d.Off)])
	}
	sum := 0.0
	for _, l := range out {
		sum += l
	}
	return out, sum
}
