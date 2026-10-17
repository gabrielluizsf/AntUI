package svg

import (
	"fmt"
	"strings"
)

// children reads what is written inside a tag until the tag that closes it.
// Text and elements are kept apart, and the whitespace around a child element
// is left out of the text, so an element that only holds formatting whitespace
// holds nothing. Each piece of writing is also kept where it was written,
// because the writing of a `<text>` runs on across the elements between it.
// preserve says the writing is kept exactly as the file wrote it, whitespace
// runs and line breaks and all, which is what `xml:space="preserve"` asks for.
func (p *xmlParser) children(e *element, preserve bool, ns namespace) error {
	var text strings.Builder
	keep := func() { keepText(e, &text, preserve) }
	for {
		if p.eof() {
			return fmt.Errorf("antui/svg: <%s> is never closed", e.Name)
		}
		if p.at("</") {
			if err := p.closeTag(e, ns); err != nil {
				return err
			}
			keep()
			e.Text = textOf(e)
			return nil
		}
		if p.at("<") {
			keep()
			kid, err := p.inside(preserve, &text, ns)
			if err != nil {
				return err
			}
			// An element in a foreign namespace is a drawing within the
			// drawing: it is read so that what follows it is read right, and
			// left out of the tree so that nothing paints it.
			if kid != nil && !kid.foreign {
				e.Kids = append(e.Kids, kid)
				e.Parts = append(e.Parts, part{Elem: kid})
			}
			continue
		}
		// The next piece of writing, an entity read into the character it
		// stands for. A name nothing defines is an error, since a drawing that
		// holds one does not say what it means.
		if err := p.writeText(e, &text); err != nil {
			return err
		}
	}
}
