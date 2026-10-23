package svg

import "fmt"

// splicePE reads a name the doctype leans on itself — written %name; — and
// reads what it holds in the place it was written, since that text is doctype
// text and may hold declarations of its own. The parser is put at the held
// text and put back where it was, and how far that text was read is no part of
// where the file itself stands.
func (p *xmlParser) splicePE(to int) error {
	word, size, err := p.peAt(p.src, p.pos, nil)
	if err != nil {
		return err
	}
	p.pos += size
	was, at := p.src, p.pos
	p.src, p.pos = word, 0
	err = p.parseSubset(0, len(word))
	p.src, p.pos = was, at
	return err
}

// peAt reads the name at the % of s[i] — a name the doctype leans on itself —
// and answers what it holds and how many bytes the whole %…; took. A name
// nothing defines, one whose text another document holds, and one that comes
// back around to itself are each an error rather than the letters they were
// written as.
func (p *xmlParser) peAt(s string, i int, over []string) (string, int, error) {
	name, size, err := refName(s, i)
	if err != nil {
		return "", 0, fmt.Errorf("antui/svg: %w", err)
	}
	e, ok := p.pes[name]
	if !ok {
		return "", 0, fmt.Errorf("antui/svg: the name %%%s; is not one this doctype defines", name)
	}
	if e.external {
		return "", 0, fmt.Errorf("antui/svg: %%%s; holds text of another document, which is not read here", name)
	}
	for _, under := range over {
		if under == name {
			return "", 0, fmt.Errorf("antui/svg: the name %%%s; comes back around to itself", name)
		}
	}
	text, err := p.expandPEOver(e.text, append(over, name))
	if err != nil {
		return "", 0, err
	}
	if err := p.spendExpansion(len(text)); err != nil {
		return "", 0, fmt.Errorf("antui/svg: %w", err)
	}
	return text, size, nil
}
