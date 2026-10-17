package svg

import "fmt"

// definedEntity reads what a name the drawing's own doctype defines holds. A
// name nothing defines, one whose text another document holds, and one whose
// text brings its own name back around are each an error: in every one of them
// there is no text here to read, and guessing would paint something other than
// the file says.
func (p *xmlParser) definedEntity(name string, size int, over []string) (string, int, error) {
	e, ok := p.ents[name]
	if !ok {
		return "", 0, fmt.Errorf("the entity &%s; is not one this drawing defines", name)
	}
	if e.external {
		return "", 0, fmt.Errorf("&%s; holds text of another document, which is not read here", name)
	}
	for _, under := range over {
		if under == name {
			return "", 0, fmt.Errorf("the entity &%s; comes back around to itself", name)
		}
	}
	text, err := p.expand(e.text, append(over, name))
	if err != nil {
		return "", 0, err
	}
	return text, size, nil
}
