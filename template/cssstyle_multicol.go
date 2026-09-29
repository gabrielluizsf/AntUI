package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintColumnRules draws the rule between two columns, once per gutter the
// container made and as tall as the tallest column. The rule is the same three
// declarations a border is, so it is drawn the same way — as side 3, a left
// edge, standing in the middle of the gutter: a width, a style (solid, dashed,
// dotted, double) and a colour, which is currentColor when the stylesheet left
// it that way.
func (cs *cssStyle) paintColumnRules(win *antui.Window, b *multiBatch) {
	st := b.st
	bw := st.ColumnRuleWidth * cs.u()
	if bw < 1 || st.ColumnRuleStyle == css.BorderNone || st.ColumnRuleColor.A() == 0 {
		return
	}
	cv := win.Canvas()
	// The rule runs the height of the columns: the tallest of them for a
	// balancing container, the whole content box for one that declared a height
	// and spends it column by column.
	top, height := b.originY, max(b.colH, b.contentH)
	for col := 1; col < b.count; col++ {
		// The gutter in front of this column is every column before it, plus
		// the gutter between each pair of them: a rule stands in the gap, not
		// inside a column.
		gutter := b.originX + col*b.colW + (col-1)*b.gap
		// Centred in the gutter; a rule wider than the gutter starts at its edge
		// and runs over the column beside it, as it does in a browser.
		x := max(gutter+(b.gap-bw)/2, gutter)
		cs.paintSide(cv, 3, st.ColumnRuleStyle, bw, x, top, bw, height, st.ColumnRuleColor)
	}
}
