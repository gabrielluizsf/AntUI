package svg

// element is one tag of a document: what it is called, what it was given, what
// was written inside it, and what was written inside those. A drawing is a tree
// of these, and the root of it is the element that has no one above.
type element struct {
	// Name is the tag with any namespace taken off, so a document that spells
	// its tags `svg:circle` and one that spells them `circle` read the same.
	Name string
	// Attr holds what the tag was given, under the names it was written with:
	// a name is read exactly as the file spells it, as XML reads one.
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
	// sheet is what the drawing's own `<style>` sheets say about this element,
	// left here by the pass that matches every rule against every element
	// before anything is built. It is nil where no rule matched, which is the
	// whole drawing that wrote no sheet.
	sheet *sheetDecls
	// foreign says the element stands in a namespace other than the SVG one,
	// so it is a drawing within the drawing and not this package's to paint.
	foreign bool
}
