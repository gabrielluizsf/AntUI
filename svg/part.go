package svg

// part is one piece of what was written inside a tag: either the writing
// between two elements, or an element. Exactly one of the two is there.
type part struct {
	Text string
	Elem *element
}
