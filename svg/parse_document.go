package svg

// parseDocument reads a whole XML document and answers its root element. The
// text is decoded first — a byte-order mark taken off, UTF-16 turned into
// UTF-8 — because what follows is read as characters, not as the bytes a file
// happened to be written in. An error is only ever structural: a tag that
// never closes, a string left open, a drawing that is not the single-rooted
// document it declared — because a drawing that does not close what it opens
// cannot be read at all, while anything else inside it is either drawn or
// reported as a warning.
func parseDocument(src string) (*element, error) {
	p := &xmlParser{src: decodeText(src)}
	return p.document()
}
