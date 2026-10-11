package template

import "github.com/gabrielluizsf/antui/template/css"

// cellPlacement is the box a cell takes in its row: the whole row's height for
// a cell that asks for the top, and its own height placed in the middle or at
// the bottom for one that asks for that.
func cellPlacement(cell *tableCell, rowHeight, rowY int) (h, y int) {
	switch cell.st.VerticalAlign {
	case css.VerticalAlignMiddle, css.VerticalAlignSub, css.VerticalAlignSuper:
		return cell.nh, rowY + (rowHeight-cell.nh)/2
	case css.VerticalAlignBottom, css.VerticalAlignTextBottom:
		return cell.nh, rowY + max(rowHeight-cell.nh, 0)
	default:
		return rowHeight, rowY
	}
}
