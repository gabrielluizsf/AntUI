package svg

// relCandidates is every node a relative selector's rightmost compound could
// stand on: the element's descendants, and the siblings written after it with
// their descendants. The chain walk then rejects the ones the lead combinator
// does not open.
func relCandidates(scope *selNode) []*selNode {
	var out []*selNode
	appendDescendants(&out, scope)
	for s := scope.next; s != nil; s = s.next {
		out = append(out, s)
		appendDescendants(&out, s)
	}
	return out
}

// appendDescendants adds every element below n to out.
func appendDescendants(out *[]*selNode, n *selNode) {
	for _, k := range n.kids {
		*out = append(*out, k)
		appendDescendants(out, k)
	}
}
