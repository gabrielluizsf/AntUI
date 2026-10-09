package css

// FlexDirection names the main axis of a flex container: which way its items
// pack. Row lays them left to right, column top to bottom, and the reverse
// variants pack from the far end of that axis.
const (
	FlexDirectionRow uint8 = iota
	FlexDirectionRowReverse
	FlexDirectionColumn
	FlexDirectionColumnReverse
)

// FlexWrap names how flex items give the main axis if they overflow it:
// nowrap lets them shrink and stick to one line, wrap grants lines, and
// wrap-reverse stacks the wrapped lines from the far end of the cross axis.
const (
	FlexWrapNowrap uint8 = iota
	FlexWrapWrap
	FlexWrapWrapReverse
)

// JustifyContent names how leftover main-axis space is spread between items:
// packed at one end, centred, or turned into spacing between, around or
// even around the items. flex-start is also the initial value.
const (
	JustifyFlexStart uint8 = iota
	JustifyFlexEnd
	JustifyCenter
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// Align* names how an item is placed across the cross axis of its line.
// AlignAuto is only legal for align-self, where it means "use the container's
// align-items"; AlignStretch is align-items' initial value, which grows an
// auto-sized item to fill its line's cross size. AlignBaseline is accepted
// and drawn as flex-start: the engine has no shared baseline model.
const (
	AlignAuto    uint8 = 0xFF
	AlignStretch uint8 = iota
	AlignFlexStart
	AlignFlexEnd
	AlignCenter
	AlignBaseline
)

// Content* name how leftover cross-axis space is distributed between whole
// lines when a container wraps. Stretch (the initial value) grows the lines,
// the rest spread or pack them.
const (
	ContentStretch uint8 = iota
	ContentFlexStart
	ContentFlexEnd
	ContentCenter
	ContentSpaceBetween
	ContentSpaceAround
	ContentSpaceEvenly
)
