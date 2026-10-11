package template

import "github.com/gabrielluizsf/antui/template/css"

// multiColumnRoom is how much content height one column may hold. A container
// that fills a column before starting the next spends its whole height on each
// one; a balancing container spreads the height over them, which is what a
// column-count with a declared height does.
func (c *CSS) multiColumnRoom(b *multiBatch, st css.Style) int {
	room := max(b.contentH, 0)
	if st.ColumnFill == css.ColumnFillAuto || b.count < 2 {
		return room
	}
	return max(room/b.count, 0)
}
