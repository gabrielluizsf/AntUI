package svg

import "fmt"

// doctype reads a `<!DOCTYPE …>`: its root name and any external id are read
// past, and its internal subset — the part between brackets — is read for the
// entities a drawing defines there, since its text leans on them. The subset
// of a DTD that lives in another document is not fetched, the same as no
// renderer fetches it, so a name only such a DTD defines stays undefined.
// Anything other than a doctype where a declaration sits does not open a
// document.
func (p *xmlParser) doctype() error {
	p.pos += len("<!")
	p.skipSpace()
	if !p.atWord("DOCTYPE") {
		return fmt.Errorf("antui/svg: a declaration here may only be a doctype, not %q", p.peek(12))
	}
	p.pos += len("DOCTYPE")
	from := -1
	for {
		p.skipSpace()
		if p.eof() {
			return fmt.Errorf("antui/svg: the doctype is never closed")
		}
		done, next, err := p.doctypeStep(from)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
		from = next
	}
}
