package template

// Hydra repaints the bound stylesheet in the Hydra skin: the near-black
// glass, the single blue accent and the flat frames the Solo Task site's HUD
// wears, with the site's own faces' type. Call it after [CSS.SetStyle],
// which reloads the sheet and would drop a skin applied before it:
//
//	app := template.TemplateWithCSS(win)
//	if err := app.SetStyle("app.css", css.NewTable()); err != nil { ... }
//	if err := app.Hydra(); err != nil { ... }
func (c *CSS) Hydra() error {
	if c == nil || c.style == nil {
		return nil
	}
	return c.style.classes.Hydra()
}
