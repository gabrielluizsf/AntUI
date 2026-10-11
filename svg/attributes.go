package svg

import (
	"fmt"
	"strings"
)

// attributes reads everything up to the closing bracket of the tag. A value may
// be in single quotes, double quotes, or no quotes at all, as XML allows the
// last when it has no space in it.
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
			// A slash with nothing before it is a tag that closes itself,
			// written the older way as `<rect/>` with no space in between.
			if p.at("/>") {
				p.pos += 2
				e.Attr[""] = "/"
				return nil
			}
			return fmt.Errorf("antui/svg: <%s> has something in it that is not an attribute", e.Name)
		}
		p.skipSpace()
		if !p.at("=") {
			// An attribute written on its own has the value its own name, which
			// is how `hidden` and `required` are written in a drawing.
			e.Attr[strings.ToLower(name)] = name
			continue
		}
		p.pos++
		p.skipSpace()
		v, err := p.attrValue()
		if err != nil {
			return fmt.Errorf("antui/svg: <%s %s>: %w", e.Name, name, err)
		}
		e.Attr[strings.ToLower(name)] = v
	}
}
