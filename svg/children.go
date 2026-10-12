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
func (p *xmlParser) children(e *element, preserve bool) error {
	var text strings.Builder
	keep := func() { keepText(e, &text, preserve) }
	for {
		if p.eof() {
			return fmt.Errorf("antui/svg: <%s> is never closed", e.Name)
		}
		if p.at("</") {
			if err := p.closeTag(e); err != nil {
				return err
			}
			keep()
			e.Text = textOf(e)
			return nil
		}
		if p.at("<!--") {
			p.through("-->")
			continue
		}
		if p.at("<![CDATA[") {
			p.readCdata(&text)
			continue
		}
		if p.at("<?") {
			p.through("?>")
			continue
		}
		if p.at("<!") {
			p.skipDecl()
			continue
		}
		if p.at("<") {
			keep()
			kid, err := p.element(preserve)
			if err != nil {
				return err
			}
			e.Kids = append(e.Kids, kid)
			e.Parts = append(e.Parts, part{Elem: kid})
			continue
		}
		text.WriteByte(p.src[p.pos])
		p.pos++
	}
}
