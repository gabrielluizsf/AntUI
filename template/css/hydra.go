package css

// Hydra repaints the loaded stylesheet in the Hydra skin: the near-black
// glass, the single blue accent and the flat frames the Solo Task site's
// HUD wears. Everything the file decided about layout is kept — the skin
// only carries colour, borders, shadows and the accent's type weight.
//
// Call it after [CSSClasses.SetStyle]:
//
//	classes := css.NewTable()
//	if err := classes.SetStyle("app.css"); err != nil { ... }
//	if err := classes.Hydra(); err != nil { ... }
func (c *CSSClasses) Hydra() error {
	if c == nil || c.sheet == nil {
		return errNoSheet
	}
	return c.Apply(hydraSkin)
}
