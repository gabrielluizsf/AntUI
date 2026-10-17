package svg

import "fmt"

// document reads from the opening angle bracket to the end of the root element
// and then to the end of the text, because a drawing that carries on after its
// root with another tag or a stray word is not the single-rooted document XML
// describes and the tree it would leave behind is not the tree that was
// declared.
func (p *xmlParser) document() (*element, error) {
	p.stripBOM()
	for {
		p.skipSpace()
		if p.at("<?") {
			p.through("?>")
			continue
		}
		if p.at("<!--") {
			p.through("-->")
			continue
		}
		if p.at("<!") {
			// A doctype or another declaration: it ends at the first > that is
			// not inside a pair of brackets, which is where the internal subset
			// of a DTD would otherwise run past it.
			p.skipDecl()
			continue
		}
		break
	}
	if p.eof() {
		return nil, fmt.Errorf("antui/svg: the drawing is empty, there is no <svg> in it")
	}
	if !p.at("<") {
		return nil, fmt.Errorf("antui/svg: the drawing starts with %q, where a tag was expected", p.peek(12))
	}
	root, err := p.element(false)
	if err != nil {
		return nil, err
	}
	if err := p.tail(); err != nil {
		return nil, err
	}
	return root, nil
}
