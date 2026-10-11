package template

// myMain reads an item's solved main-axis size.
func (b *flexBatch) myMain(it *flexItem) int {
	if b.col() {
		return it.h
	}
	return it.w
}
