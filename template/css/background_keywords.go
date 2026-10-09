package css

// The gradient families a background-image may hold.
const (
	GradientLinear uint8 = iota
	GradientRadial
	GradientConic
)

// The radial ending shapes.
const (
	GradientEllipse uint8 = iota
	GradientCircle
)

// The radial size keywords, from the CSS farthest-corner default.
const (
	GradientFarthestCorner uint8 = iota
	GradientClosestSide
	GradientFarthestSide
	GradientClosestCorner
)

// BackBox names one member of the border-box / padding-box / content-box
// family, the box a background-origin anchors to or a background-clip cuts
// the paint at.
const (
	BackBorder uint8 = iota
	BackPadding
	BackContent
)

// The background-repeat keywords for one axis.
const (
	BackRepeatRepeat uint8 = iota
	BackRepeatNoRepeat
	BackRepeatSpace
	BackRepeatRound
)

// The background-attachment keywords; the engine only models them, painting
// every layer at its box as if it were scroll.
const (
	BackAttachScroll uint8 = iota
	BackAttachFixed
	BackAttachLocal
)
