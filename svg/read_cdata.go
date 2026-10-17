package svg

import (
	"fmt"
	"strings"
)

// readCdata appends the text of a `<![CDATA[ … ]]>` section to b and leaves
// the position past it. The section ends where `]]>` begins, so the text is
// what is between the two rather than up to the end of them. A section that
// was never closed is an error rather than the rest of the file quietly
// becoming text.
func (p *xmlParser) readCdata(b *strings.Builder) error {
	p.pos += len("<![CDATA[")
	i := strings.Index(p.src[p.pos:], "]]>")
	if i < 0 {
		return fmt.Errorf("antui/svg: a CDATA section is never closed")
	}
	b.WriteString(p.src[p.pos : p.pos+i])
	p.pos += i + len("]]>")
	return nil
}
