package svg

import "fmt"

// element reads one tag and everything inside it, up to its own closing tag.
// preserve is what the element above said about its own writing, carried down
// because `xml:space` is inherited; the element itself may say otherwise. ns
// is what the prefixes in scope stand for, and the element's own bindings are
// in scope for its own name: a tag may bind the prefix it is written with.
func (p *xmlParser) element(preserve bool, ns namespace) (*element, error) {
	p.through("<")
	start := p.pos
	for !p.eof() && !isSpace(p.src[p.pos]) && p.src[p.pos] != '>' && p.src[p.pos] != '/' {
		p.pos++
	}
	name := p.src[start:p.pos]
	if name == "" {
		return nil, fmt.Errorf("antui/svg: a tag has no name at offset %d", start)
	}
	if !validQName(name) {
		return nil, fmt.Errorf("antui/svg: a tag is written as %q at offset %d, which is not a name XML allows", name, start)
	}
	e := &element{Name: name, Attr: map[string]string{}}
	if err := p.attributes(e); err != nil {
		return nil, err
	}
	scope := ns.child(e)
	e.Name, e.foreign = scope.resolve(name)
	if e.Attr[""] != "" {
		// A lone slash before the closing bracket is the older way of writing a
		// tag that closes itself, and it means the same thing.
		return e, nil
	}
	if err := p.children(e, keepsSpace(e, preserve), scope); err != nil {
		return nil, err
	}
	return e, nil
}
