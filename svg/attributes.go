package svg

import (
	"fmt"
	"strings"
)

// attributes reads everything up to the closing bracket of the tag. A value may
// be quoted or not — a drawing written by hand does the last when it holds no
// space, which XML does not ask for but SVG in the wild arrives as. A name
// repeated on one tag is an error rather than the second value quietly winning,
// and a name that is not a name at all is the same sort of error.
func (p *xmlParser) attributes(e *element) error {
	for {
		p.skipSpace()
		if p.eof() {
			return fmt.Errorf("antui/svg: <%s> is never closed", e.Name)
		}
		if p.at("/>") {
			p.pos += 2
			e.Attr[""] = "/"
			return nil
		}
		if p.at(">") {
			p.pos++
			return nil
		}
		name, ok := p.attrName()
		if !ok {
			// A slash with nothing before it closes the tag the older way.
			if p.at("/>") {
				p.pos += 2
				e.Attr[""] = "/"
				return nil
			}
			return fmt.Errorf("antui/svg: <%s> has something in it that is not an attribute", e.Name)
		}
		if !validName(name) {
			return fmt.Errorf("antui/svg: <%s> has an attribute %q, which is not a name XML allows", e.Name, name)
		}
		key := strings.ToLower(name)
		if _, seen := e.Attr[key]; seen && key != "" {
			return fmt.Errorf("antui/svg: <%s> is given %s twice", e.Name, name)
		}
		p.skipSpace()
		if !p.at("=") {
			// An attribute written on its own means what its own name means.
			e.Attr[key] = name
			continue
		}
		p.pos++
		p.skipSpace()
		v, err := p.attrValue()
		if err != nil {
			return fmt.Errorf("antui/svg: <%s %s>: %w", e.Name, name, err)
		}
		e.Attr[key] = v
	}
}
