package template

import "github.com/gabrielluizsf/antui/canvas"

// Label draws a line of text in the template's own voice.
func (c *Context) Label(x, y int, text string) { c.tpl.Label(c.win, canvas.Area{X: x, Y: y}, text) }
