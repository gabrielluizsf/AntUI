package svg

// find is the first node of that name anywhere at or below this one, and nil
// where there is none. It is how a caller reaches one shape of a drawing to
// measure it or to see what style the file gave it.
func (n *Node) find(name string) *Node {
	if n == nil {
		return nil
	}
	if n.Name == name {
		return n
	}
	for _, k := range n.Kids {
		if got := k.find(name); got != nil {
			return got
		}
	}
	return nil
}
