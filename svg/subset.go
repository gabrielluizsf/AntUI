package svg

import "fmt"

// parseSubset reads the declarations between the doctype's brackets: the
// entities a drawing defines there are kept, since its text leans on them, and
// every other declaration is read past, being a matter of validity this
// renderer does not keep to. A name written %name; here holds pieces of
// doctype text and is read in the place it was written.
func (p *xmlParser) parseSubset(from, to int) error {
	p.pos = from
	for p.pos < to {
		p.skipSpace()
		if p.pos >= to {
			break
		}
		if p.at("<!--") {
			if err := p.comment(); err != nil {
				return err
			}
			continue
		}
		if p.at("<?") {
			if err := p.pi(); err != nil {
				return err
			}
			continue
		}
		if p.at("<!") {
			if err := p.subsetDecl(to); err != nil {
				return err
			}
			continue
		}
		if p.at("%") {
			if err := p.splicePE(to); err != nil {
				return err
			}
			continue
		}
		return fmt.Errorf("antui/svg: %q is written where a declaration of the doctype was expected", p.peek(12))
	}
	p.pos = to
	return nil
}
