package svg

// selTree builds the view a selector walks, in the order the elements were
// written, each knowing the one above it and the ones beside it.
func selTree(e *element, parent, prev *selNode) *selNode {
	n := &selNode{el: e, parent: parent, prev: prev}
	if prev != nil {
		prev.next = n
	}
	var cp *selNode
	for _, k := range e.Kids {
		child := selTree(k, n, cp)
		n.kids = append(n.kids, child)
		cp = child
	}
	return n
}
