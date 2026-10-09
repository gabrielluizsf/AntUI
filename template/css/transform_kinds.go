package css

// The kinds a [TransformFunc] can be. rotate() and skew() carry their angles
// in the Ax/Ay fields; a matrix() fills M; scale() lands in Sx/Sy and
// translate() keeps Lengths so percentages stay percentages until the box is
// known.
const (
	TransformTranslate uint8 = iota
	TransformScale
	TransformRotate
	TransformSkew
	TransformMatrix
)
