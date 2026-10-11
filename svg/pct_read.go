package svg

// pctRead is one number written as a percentage, kept until the units it is a
// percentage of are known: what was written, and which axis of the drawing it is
// measured along. The axis is -1 for the one a percentage has no side to it,
// which is a radius.
type pctRead struct {
	raw  string
	axis int
}
