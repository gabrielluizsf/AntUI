package template

// baseMain picks an item's base main-axis size: its flex-basis when the basis
// is definite, its natural size otherwise. A percentage basis needs a
// definite axis to resolve against, or it falls back to the natural size.
func (b *flexBatch) baseMain(it *flexItem, availMain int, definite bool) int {
	natural := it.nw
	if b.col() {
		natural = it.nh
	}
	base := natural
	if !it.basis.Auto() && !it.basis.None() {
		if !(it.basis.IsPct() && !definite) {
			base = b.css.style.length(b.st, it.basis, availMain)
		}
	}
	if it.maxMain >= 0 {
		base = min(base, it.maxMain)
	}
	if it.minMain >= 0 {
		base = max(base, it.minMain)
	}
	return base
}
