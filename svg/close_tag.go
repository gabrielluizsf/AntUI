package svg

import "fmt"

// closeTag reads the `</name>` that closes e and checks that it is e's own
// name: a drawing that opens a group and closes it with something else has no
// shape to paint, and reading it as though it were right would draw something
// other than what the file says.
func (p *xmlParser) closeTag(e *element) error {
	p.pos += 2
	start := p.pos
	for !p.eof() && p.src[p.pos] != '>' {
		p.pos++
	}
	if p.eof() {
		return fmt.Errorf("antui/svg: <%s> is never closed", e.Name)
	}
	if closing := localName(p.src[start:p.pos]); closing != e.Name {
		return fmt.Errorf("antui/svg: <%s> is closed by </%s>", e.Name, closing)
	}
	p.pos++
	return nil
}
