package svg

// lengthUnits are the units a length may be written in, as how many pixels one
// of them is. A unitless number is not in here: it is the drawing's own, and
// only the caller knows how big that is.
//
// The longest unit comes first, because one of them ends in another: `cm` and
// `mm` share their last letter, and `in` is the tail of nothing else here, but a
// number is read by the first unit it genuinely ends in rather than the first
// whose last letter happens to match.
var lengthUnits = []struct {
	suffix string
	scale  float64
}{
	{"px", 1},
	{"pt", 96.0 / 72},
	{"pc", 16},
	{"cm", 96 / 2.54},
	{"mm", 96 / 25.4},
	{"in", 96},
	{"q", 96 / 101.6},
}
