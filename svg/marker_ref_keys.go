package svg

// markerRefXKeys and markerRefYKeys are the geometric keywords refX and refY
// accept, as the fraction of the length they name that each stands for. SVG 2
// spells the reference point with the same words that pin a picture to its box;
// a keyword belonging to the other axis is not in this table, and is said like
// any other value that is not a length, which is what an invalid document
// deserves.
var (
	markerRefXKeys = map[string]float64{"left": 0, "center": 0.5, "right": 1}
	markerRefYKeys = map[string]float64{"top": 0, "center": 0.5, "bottom": 1}
)
