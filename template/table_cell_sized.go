package template

import "github.com/gabrielluizsf/antui/template/css"

// tableCellSized is a cell whose content is a container: it asks for the box
// the container asked for, and the row gives it the width of its column.
func (c *CSS) tableCellSized(st css.Style, nw, nh int) *tableCell {
	return &tableCell{st: st, nw: nw, nh: nh, w: nw, h: nh}
}
