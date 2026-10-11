package template

// pieceAt is the fragment the draw pass is on for one item: the piece of the
// pass's own number, and the item's last one once the flow has nothing further
// to cut.
func (b *multiBatch) pieceAt(it *multiItem) multiPiece {
	p := min(b.pass, len(it.pieces)-1)
	if p < 0 {
		p = 0
	}
	return it.pieces[p]
}
