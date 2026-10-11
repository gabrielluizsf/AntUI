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
	keep := func() {
		s := text.String()
		if !preserve {
			s = collapseSpace(s)
		}
		if s != "" {
			e.Parts = append(e.Parts, part{Text: s})
		}
		text.Reset()
	}
	for {
		if p.eof() {
			return fmt.Errorf("antui/svg: <%s> is never closed", e.Name)
		}
		if p.at("</") {
			p.pos += 2
			start := p.pos
			for !p.eof() && p.src[p.pos] != '>' {
				p.pos++
			}
			if p.eof() {
				return fmt.Errorf("antui/svg: <%s> is never closed", e.Name)
			}
			// The tag that closes this one has to be the one this one is: a
			// drawing that opens a group and closes it with something else has
			// no shape to paint, and reading it as though it were right would
			// draw something other than what the file says.
			if closing := localName(p.src[start:p.pos]); closing != e.Name {
				return fmt.Errorf("antui/svg: <%s> is closed by </%s>", e.Name, closing)
			}
			p.pos++
			keep()
			e.Text = textOf(e)
			return nil
		}
		if p.at("<!--") {
			p.through("-->")
			continue
		}
		if p.at("<![CDATA[") {
			p.pos += len("<![CDATA[")
			start := p.pos
			// The section ends where `]]>` begins, so the text is what is
			// between the two rather than up to the end of them.
			if i := strings.Index(p.src[p.pos:], "]]>"); i >= 0 {
				text.WriteString(p.src[start : start+i])
				p.pos += i + len("]]>")
			} else {
				text.WriteString(p.src[start:])
				p.pos = len(p.src)
			}
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
