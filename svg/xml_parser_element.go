package svg

import "fmt"

// element reads one tag and everything inside it, up to its own closing tag.
// preserve is what the element above said about its own writing, carried down
// because `xml:space` is inherited; the element itself may say otherwise.
func (p *xmlParser) element(preserve bool) (*element, error) {
	p.through("<")
	start := p.pos
	for !p.eof() && !isSpace(p.src[p.pos]) && p.src[p.pos] != '>' && p.src[p.pos] != '/' {
		p.pos++
	}
	name := p.src[start:p.pos]
	if name == "" {
		return nil, fmt.Errorf("antui/svg: a tag has no name at offset %d", start)
	}
	e := &element{Name: localName(name), Attr: map[string]string{}}
	if err := p.attributes(e); err != nil {
		return nil, err
	}
	if e.Attr[""] != "" {
		// A lone slash before the closing bracket is the older way of writing a
		// tag that closes itself, and it means the same thing.
		return e, nil
	}
	if err := p.children(e, keepsSpace(e, preserve)); err != nil {
		return nil, err
	}
	return e, nil
}
