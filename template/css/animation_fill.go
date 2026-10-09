package css

// The fill-mode of an animation: whether the frames bleed over the delay and
// the end, so a button that animates once keeps its final colour. The zero
// is none.
const (
	FillNone uint8 = iota
	FillBackwards
	FillForwards
	FillBoth
)
