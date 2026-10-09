package css

// TransformOrigin puts the pivot of a transform in the border box. The two
// lengths are the x and y of the pivot, a percentage of the box; the initial
// value is "50% 50%", the centre.
type TransformOrigin = [2]Length

const TransformOriginPlain = "transform-origin"

var InitialTransformOrigin = [2]Length{Pct(50), Pct(50)}
