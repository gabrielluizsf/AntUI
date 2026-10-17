package svg

import (
	"fmt"
	"strings"
)

// inside reads one piece of what stands inside a tag: a comment, a CDATA
// section or an instruction, each passed over with nothing to show for it; a
// declaration, which is not allowed there at all; or a tag, whose element
// comes back as the child it is. The text builder takes what a CDATA section
// writes.
func (p *xmlParser) inside(preserve bool, text *strings.Builder) (*element, error) {
	switch {
	case p.at("<!--"):
		return nil, p.comment()
	case p.at("<![CDATA["):
		return nil, p.readCdata(text)
	case p.at("<?"):
		if p.atXMLDecl() {
			return nil, fmt.Errorf("antui/svg: the xml declaration may only open the document, not sit within it")
		}
		return nil, p.pi()
	case p.at("<!"):
		return nil, fmt.Errorf("antui/svg: a declaration may not stand inside the drawing, where only a comment, a CDATA section or a tag may")
	}
	return p.element(preserve)
}
