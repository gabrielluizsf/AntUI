package css

// gridStyle carries the grid fields of a Style: the tracks it is divided
// into, the areas it names and how each item places itself.
type gridStyle struct {
	GridTemplateColumns GridTemplate
	GridTemplateRows    GridTemplate
	GridTemplateAreas   [][]string
	GridAutoFlow        uint8
	GridAutoFlowDense   bool
	GridJustifyContent  uint8
	JustifyItems        uint8
	JustifySelf         uint8
	GridColumn          GridPlacement
	GridRow             GridPlacement
}
