package css

import "github.com/gabrielluizsf/antui/canvas"

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

// BackPosition is one axis of a background-position, or a gradient's centre:
// an edge to anchor the picture to and an offset away from it. The offset is
// a length in reference pixels or a percentage of the side left over — the
// positioning area minus the image, as CSS measures it.
type BackPosition struct {
	Edge uint8 // one of the BackPos* constants
	Off  Length
}

// The edges a BackPosition may anchor to. BackPosOffset is a bare length or
// percentage with no keyword, measured from the starting edge like the
// BackPosStart default.
const (
	BackPosStart  uint8 = iota // the left/top edge
	BackPosMiddle              // the centre
	BackPosEnd                 // the right/bottom edge
	BackPosOffset              // a bare length/percentage, no keyword
)

// BackPos is a background-position: the picture's horizontal and vertical
// anchors. The zero value — both axes anchored to their start edges with no
// offset — is the CSS initial of 0% 0%.
type BackPos [2]BackPosition

// BackSize is a background-size: explicit widths in each axis, or the cover
// and contain keywords. The zero value is "auto auto", the picture at its own
// size.
type BackSize struct {
	W, H    Length
	Cover   bool
	Contain bool
}

// BackRepeat is the per-axis repeat behaviour of one background layer.
type BackRepeat [2]uint8 // X, Y; the BackRepeat* constants

// Stop is one colour stop of a gradient. Offset is a fraction along the
// gradient line, or -1 when the declaration left it out and the painter is to
// spread it; Color may be [CurrentColor] until the cascade resolves it.
type Stop struct {
	Offset float64
	Color  canvas.Color
}

// Gradient is one parsed gradient function. Angle is radians measured
// clockwise from the top, matching "to top" being 0°; Center anchors the
// radial or conic centre to a point in its painting area, with the default
// (built by the parsers) the middle. Shape and Size are radial-only.
type Gradient struct {
	Kind   uint8 // one of the Gradient* constants
	Angle  float64
	Center BackPos
	Shape  uint8
	Size   uint8
	Stops  []Stop
}

// BackImage is one background layer: either a url() the template looks up in
// its image registry, or a gradient painted directly.
type BackImage struct {
	URL  string
	Grad *Gradient
}

// backMiddle is the default gradient centre: the middle of each axis, with no
// offset.
func backMiddle() BackPos {
	return BackPos{
		{Edge: BackPosMiddle, Off: Zero()},
		{Edge: BackPosMiddle, Off: Zero()},
	}
}
