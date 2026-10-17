package svg

import "fmt"

// subsetDecl reads one declaration inside the doctype: the ENTITY one a
// drawing leans on is read and kept, and the rest belongs to validation this
// renderer does not do and is read past.
func (p *xmlParser) subsetDecl(to int) error {
	p.pos += len("<!")
	if p.atWord("ENTITY") {
		return p.entityDecl(to)
	}
	return p.skipDecl(to)
}

// skipDecl reads past a declaration this renderer keeps nothing of, where a
// `>` inside a quoted part is a character and not the end of it.
func (p *xmlParser) skipDecl(to int) error {
	for p.pos < to {
		switch p.src[p.pos] {
		case '>':
			p.pos++
			return nil
		case '"', '\'':
			if err := p.inQuotes(); err != nil {
				return err
			}
		default:
			p.pos++
		}
	}
	return fmt.Errorf("antui/svg: a declaration of the doctype is never closed")
}
