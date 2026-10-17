package svg

import (
	"fmt"
	"strings"
)

// writeText reads one piece of the writing inside a tag into text, an entity
// being read as the character it stands for.
func (p *xmlParser) writeText(e *element, text *strings.Builder) error {
	if p.src[p.pos] == '&' {
		word, size, err := p.entityAt(p.src, p.pos, nil)
		if err != nil {
			return fmt.Errorf("antui/svg: the writing of <%s>: %w", e.Name, err)
		}
		text.WriteString(word)
		p.pos += size
		return nil
	}
	// `]]>` may only close a CDATA section, never stand in a drawing's own
	// writing, where XML keeps it out.
	if p.src[p.pos] == ']' && p.at("]]>") {
		return fmt.Errorf("antui/svg: the writing of <%s> holds a `]]>`, which XML keeps out of text", e.Name)
	}
	text.WriteByte(p.src[p.pos])
	p.pos++
	return nil
}
