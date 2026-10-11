package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// newCSSStyle makes a cssStyle over an empty class table; the table is shared
// with its template, so SetStyle on the template reaches the painter.
func newCSSStyle(win *antui.Window) *cssStyle {
	return &cssStyle{
		win:      win,
		classes:  css.NewTable(),
		images:   map[string]*canvas.Canvas{},
		timeline: map[string]*widgetEntry{},
		now:      antui.Now,
	}
}
