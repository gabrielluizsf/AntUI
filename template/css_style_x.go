package template

import "github.com/gabrielluizsf/antui/template/css"

// Style exposes the underlying class table, so the developer can reach the
// whole [css.CSSClasses] API — Apply, SetProperty, Property, Warnings — after
// a stylesheet is loaded. The returned pointer is shared with the painter.
func (c *CSS) Style() *css.CSSClasses { return c.style.classes }
