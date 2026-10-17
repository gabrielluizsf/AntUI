package svg

// xmlParser walks the text of a document once, handing back what it finds. It
// is not a general XML parser: what a drawing never uses — an external DTD, an
// instruction that declares something — is read past rather than read into.
// What a drawing may lean on, it reads: the doctype's internal subset for the
// entities a drawing defines for itself, the XML declaration for the encoding
// it names, and the namespaces its tags and attributes are written against.
type xmlParser struct {
	src string
	pos int
	// Name tables from the doctype's internal subset: a drawing's own names,
	// and the doctype's parameter ones, which hold pieces of the doctype.
	ents entityTable
	pes  entityTable
}
