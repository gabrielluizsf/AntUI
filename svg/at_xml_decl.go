package svg

import "strings"

// atXMLDecl reports whether the instruction at the cursor is the xml
// declaration, which is the reserved target `xml` followed by the space or the
// close that ends a target. `<?xml-stylesheet …?>` is a plain instruction and
// not this one, and so is anything else whose name merely starts with xml.
func (p *xmlParser) atXMLDecl() bool {
	if !p.at("<?xml") {
		return false
	}
	rest := p.src[p.pos+len("<?xml"):]
	return rest == "" || isSpace(rest[0]) || strings.HasPrefix(rest, "?>")
}
