package template

// setMain stores the solved main-axis size; the cross stays where it was.
func (b *flexBatch) setMain(it *flexItem, size int) {
	if it.maxMain >= 0 {
		size = min(size, it.maxMain)
	}
	if it.minMain >= 0 {
		size = max(size, it.minMain)
	}
	if b.col() {
		it.h = size
		return
	}
	it.w = size
}
