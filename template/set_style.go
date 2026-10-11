package template

import "github.com/gabrielluizsf/antui/template/css"

// SetStyle loads the stylesheet into the table and binds it to every widget
// the template draws afterwards. It replaces whatever the table held, and any
// compile error stops the load.
//
//	classes := css.NewTable()
//	classes.Button = "primary wide"
//	if err := tpl.SetStyle("app.css", classes); err != nil { ... }
func (c *CSS) SetStyle(cssFile string, classes *css.CSSClasses) error {
	if err := classes.SetStyle(cssFile); err != nil {
		return err
	}
	c.style.classes = classes
	c.cursorY = 0
	c.lineX = 0
	c.lineH = 0
	return nil
}
