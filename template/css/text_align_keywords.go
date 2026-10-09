package css

// VerticalAlign* are the vertical-align keywords. A length is stored apart,
// in Style.BaselineShift.
const (
	VerticalAlignBaseline uint8 = iota
	VerticalAlignSub
	VerticalAlignSuper
	VerticalAlignMiddle
	VerticalAlignTop
	VerticalAlignBottom
	VerticalAlignTextTop
	VerticalAlignTextBottom
)

// The absolute weights normal and bold name; the template draws anything at
// 600 or above as bold.
const (
	FontWeightNormal uint16 = 400
	FontWeightBold   uint16 = 700
)
