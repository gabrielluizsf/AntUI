package svg

// attr reads one attribute, answering the empty string for one that was not
// there, so a caller can ask without checking twice.
//
// The name is matched exactly as it is written here, because that is how XML
// reads a name: `viewBox` and `viewbox` are two different attributes, and only
// the one the drawing is read by means anything — a file spelling one in
// another case has said nothing about it, the same as it says nothing to a
// browser.
func (e *element) attr(name string) string { return e.Attr[name] }
