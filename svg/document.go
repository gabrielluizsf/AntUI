package svg

import "fmt"

// document reads from the opening angle bracket to the end of the root element.
func (p *xmlParser) document() (*element, error) {
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
	return root, nil
}
