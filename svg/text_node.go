package svg

// textNode builds a `<text>`: the writing flattened into runs in the order it
// was written, rather than a node for every `<tspan>` in it, because a
// `<tspan>` is not a shape of its own — it is a stretch of the writing with a
// place and a style of its own, and it only means anything beside the writing
// on either side of it.
func (img *Image) textNode(e *element, st Style, warn func(string, ...any)) *Node {
	n := &Node{Name: e.Name, Style: st}
	n.Runs = trimEdges(img.textRuns(e, st, warn), st.PreserveSpace)
	// The letters are drawn one bitmap at a time from where their pen lands, so
	// a drawing that turns the writing has each letter turned with it at paint
	// time. A turn with no way back is the one that cannot be put right: the
	// matrix would have to be undone to know where a letter's pixels land, and
	// there is no such place. Saying so is better than writing somewhere the
	// file did not ask for.
	if (st.Transform.B != 0 || st.Transform.C != 0) && len(n.Runs) > 0 {
		if _, ok := st.Transform.Inverse(); !ok {
			warn("the writing is turned by a transform with no way back, so it lands where the pen is rather than along the turn")
		}
	}
	return n
}
