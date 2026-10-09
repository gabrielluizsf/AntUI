package css

// flexStyle carries the flex fields of a Style: how a DisplayFlex box lays
// its children out along the main axis and across it, and how this box
// answers its own flex container. RowGap and ColumnGap space the items and
// the wrapped lines. Order reorders items; grow and shrink share the free
// space; basis is the item's main size before distribution; align-self
// overrides the container's align-items.
type flexStyle struct {
	FlexDirection  uint8 // one of the FlexDirection* constants
	FlexWrap       uint8 // one of the FlexWrap* constants
	JustifyContent uint8 // one of the Justify* constants
	AlignItems     uint8 // one of the Align* constants
	AlignContent   uint8 // one of the Content* constants
	RowGap         Length
	ColumnGap      Length

	Order      int
	FlexGrow   float64
	FlexShrink float64
	FlexBasis  Length
	AlignSelf  uint8 // AlignAuto or one of the Align* constants
}
