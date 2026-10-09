package css

// NewTable makes an empty class table, with no stylesheet behind it. Load the
// sheet with [CSSClasses.SetStyle] before drawing anything.
func NewTable() *CSSClasses {
	return &CSSClasses{specs: make(map[styleKey]Style)}
}
