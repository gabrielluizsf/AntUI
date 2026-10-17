package svg

import (
	"fmt"
	"strings"
)

// entityDecl reads a `<!ENTITY …>`: the name, and then — in the order it is
// written — either the text that name holds or a word saying the text is
// another document's to read, which is not read here, so a drawing leaning on
// such a name leans on nothing.
func (p *xmlParser) entityDecl(to int) error {
	p.pos += len("ENTITY")
	p.skipSpace()
	parameter := false
	if p.at("%") {
		parameter = true
		p.pos++
		p.skipSpace()
	}
	name := p.readName()
	if name == "" {
		return fmt.Errorf("antui/svg: a `<!ENTITY` is never given a name")
	}
	p.skipSpace()
	if p.atWord("SYSTEM") || p.atWord("PUBLIC") {
		if err := p.skipDecl(to); err != nil {
			return err
		}
		return p.declareEntity(name, entity{external: true}, parameter)
	}
	text, err := p.entityValue(to)
	if err != nil {
		return fmt.Errorf("antui/svg: the entity &%s;: %w", name, err)
	}
	if !p.at(">") {
		return fmt.Errorf("antui/svg: the entity &%s; is never closed", name)
	}
	p.pos++
	// A `<` in the text of a drawing's own entity would be a tag wherever the
	// name is used, and a tag an entity hides is not one a reader can see. A
	// name the doctype leans on itself holds doctype text, where markup is the
	// point.
	if !parameter && strings.ContainsRune(text, '<') {
		return fmt.Errorf("antui/svg: the entity &%s; holds a `<`, which no text of a drawing may hold", name)
	}
	return p.declareEntity(name, entity{text: text}, parameter)
}
