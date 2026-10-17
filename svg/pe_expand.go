package svg

import "strings"

// expandPE reads the names a doctype leans on itself — written %name; — and
// puts what they hold in their place. A name nothing defines, or one whose
// text is another document's to read, is an error: a drawing's text leans on
// what the doctype holds, and reading past one would lean on nothing.
func (p *xmlParser) expandPE(s string) (string, error) {
	return p.expandPEOver(s, nil)
}

// expandPEOver reads the names inside a piece of doctype text, where over
// holds the names whose text is being read right now: one whose text brings
// its own name back around would never stop reading.
func (p *xmlParser) expandPEOver(s string, over []string) (string, error) {
	if !strings.ContainsRune(s, '%') {
		return s, nil
	}
	var out strings.Builder
	out.Grow(len(s))
	for i := 0; i < len(s); {
		if s[i] != '%' {
			out.WriteByte(s[i])
			i++
			continue
		}
		word, size, err := p.peAt(s, i, over)
		if err != nil {
			return "", err
		}
		out.WriteString(word)
		i += size
	}
	return out.String(), nil
}
