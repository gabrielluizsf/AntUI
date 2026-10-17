package svg

import "fmt"

// document reads from the opening angle bracket to the end of the root element
// and then to the end of the text, because a drawing that carries on after its
// root with another tag or a stray word is not the single-rooted document XML
// describes. The xml declaration may open the document and only there, as XML
// has it; a comment or an instruction that was never closed stops the reading
// instead of swallowing the file.
func (p *xmlParser) document() (*element, error) {
	if p.atXMLDecl() {
		if err := p.declaration(); err != nil {
			return nil, err
		}
	}
	for {
		p.skipSpace()
		switch {
		case p.eof():
			return nil, fmt.Errorf("antui/svg: the drawing is empty, there is no <svg> in it")
		case p.atXMLDecl():
			return nil, fmt.Errorf("antui/svg: the xml declaration may only open the document, not sit within it")
		case p.at("<!--"):
			if err := p.comment(); err != nil {
				return nil, err
			}
		case p.at("<?"):
			if err := p.pi(); err != nil {
				return nil, err
			}
		case p.at("<!"):
			// A doctype or another declaration: it ends at the first > that is
			// not inside a pair of brackets, which is where the internal subset
			// of a DTD would otherwise run past it.
			p.skipDecl()
		case p.at("<"):
			root, err := p.element(false, namespace{})
			if err != nil {
				return nil, err
			}
			if err := p.tail(); err != nil {
				return nil, err
			}
			return root, nil
		default:
			return nil, fmt.Errorf("antui/svg: the drawing starts with %q, where a tag was expected", p.peek(12))
		}
	}
}
