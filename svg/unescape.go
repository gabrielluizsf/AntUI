package svg

import "strings"

// unescape reads an attribute value as the drawing means it: the five entities
// XML defines, the characters written as `&#…;`, and any name the drawing's
// own doctype defines are each read into the place they were written. A bare
// `&`, or a name nothing defines, is an error — a drawing whose text holds one
// does not say what it means, and reading past it would paint something other
// than the file.
func (p *xmlParser) unescape(s string) (string, error) {
	return p.expand(s, nil)
}

// expand reads the entities inside a piece of text, where over holds the names
// whose text is being read right now: one whose text brings its own name back
// around would never stop reading, so it is an error instead.
func (p *xmlParser) expand(s string, over []string) (string, error) {
	if !strings.ContainsRune(s, '&') {
		return s, nil
	}
	var out strings.Builder
	out.Grow(len(s))
	for i := 0; i < len(s); {
		if s[i] != '&' {
			out.WriteByte(s[i])
			i++
			continue
		}
		word, size, err := p.entityAt(s, i, over)
		if err != nil {
			return "", err
		}
		out.WriteString(word)
		i += size
	}
	return out.String(), nil
}
