package svg

// scale2 is how far a transform stretches each axis: the length the unit x and
// the unit y directions come out at. A transform that only turns or moves
// leaves both at one; one that stretches an axis more than the other says so in
// the two numbers, which is what a blur or a shadow offset off one axis needs
// and a single average cannot say.
type scale2 struct{ x, y float64 }
