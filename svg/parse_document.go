package svg

// parseDocument reads a whole XML document and answers its root element. An
// error is only ever structural — a tag that never closes, a string left open —
// because a drawing that does not close what it opens cannot be read at all,
// while anything else inside it is either drawn or reported as a warning.
func parseDocument(src string) (*element, error) {
	p := &xmlParser{src: src}
	return p.document()
}
