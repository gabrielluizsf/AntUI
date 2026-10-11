package template

// ask is how wide a column block being measured says it wants to be to whatever
// holds it. A block fills the room it is given, but the room a measure pass
// offers is a guess — in a cell whose columns are not sized it is the whole
// window — so a width in a percentage answers with what the columns inside need
// instead: one column's worth, or the widest box when there is a single column.
// That is what a percentage with no definite box to resolve against counts as,
// and it is what sizes the column that holds the block. The draw pass solves it
// again at the width the parent really gave.
func (b *multiBatch) ask() int {
	if !b.st.Has("width") || !b.st.Width.IsPct() {
		return b.w
	}
	w := b.colW
	if b.count < 2 {
		w = 0
		for i := range b.items {
			w = max(w, b.items[i].nw)
		}
	}
	return w + b.left + b.right
}
