package svg

import "fmt"

// closeTag reads the `</name>` that closes e and checks that it is e's own
// name: a drawing that opens a group and closes it with something else has no
// shape to paint, and reading it as though it were right would draw something
// other than what the file says. The name is read through the same scope as
// the tag it closes, so a prefix resolves the same way on both sides. XML
// lets whitespace sit between the name and the bracket — `</g >` closes a
// `<g>` — so the name ends at the space, not at the bracket.
func (p *xmlParser) closeTag(e *element, ns namespace) error {
	p.pos += 2
	start := p.pos
	for !p.eof() && !isSpace(p.src[p.pos]) && p.src[p.pos] != '>' {
		p.pos++
	}
	closing, foreign, err := ns.resolve(p.src[start:p.pos])
	if err != nil {
		return err
	}
	p.skipSpace()
	if p.eof() || p.src[p.pos] != '>' {
		return fmt.Errorf("antui/svg: <%s> is closed by </%s>, which is not a closing tag", e.Name, closing)
	}
	if closing != e.Name || foreign != e.foreign {
		return fmt.Errorf("antui/svg: <%s> is closed by </%s>", e.Name, closing)
	}
	p.pos++
	return nil
}
