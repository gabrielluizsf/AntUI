package svg

// xmlParser walks the text of a document once, handing back what it finds. It
// is not a general XML parser: the parts of the language a drawing never uses,
// entities beyond the five XML defines and character references, DTDs and
// processing instructions, are read past rather than read into.
type xmlParser struct {
	src string
	pos int
}
