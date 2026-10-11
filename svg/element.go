package svg

// element is one tag of a document: what it is called, what it was given, what
// was written inside it, and what was written inside those. A drawing is a tree
// of these, and the root of it is the element that has no one above.
type element struct {
	// Name is the tag with any namespace taken off, so a document that spells
	// its tags `svg:circle` and one that spells them `circle` read the same.
	Name string
	// Attr holds what the tag was given, keys lowercased because attribute
	// names are not case sensitive.
	Attr map[string]string
	// Kids are the elements written inside this one, in the order they came.
	Kids []*element
	// Text is the writing between the tags, with the whitespace at the two ends
	// of it taken out, because whitespace around a drawing is there to make
	// it readable and says nothing about where it is. It is the whole of the
	// writing, put together; see Parts for the order it was written in.
	Text string
	// Parts are what was written inside the tag, in the order it was written:
	// the writing and the elements one after another. A `<text>` holds both,
	// and reading it as one string beside a list of children would put the
	// writing on either side of a `<tspan>` in the wrong order.
	Parts []part
}
