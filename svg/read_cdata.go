package svg

import "strings"

// readCdata appends the text of a `<![CDATA[ ... ]]>` section to b and leaves
// the position past it. The section ends where `]]>` begins, so the text is
// what is between the two rather than up to the end of them.
func (p *xmlParser) readCdata(b *strings.Builder) {
	p.pos += len("<![CDATA[")
	start := p.pos
	if i := strings.Index(p.src[p.pos:], "]]>"); i >= 0 {
		b.WriteString(p.src[start : start+i])
		p.pos += i + len("]]>")
	} else {
		b.WriteString(p.src[start:])
		p.pos = len(p.src)
	}
}
