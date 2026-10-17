package svg

import "fmt"

// attrValue reads the value of one attribute, unwrapping its quotes and turning
// what is inside them into the characters it stands for. A `<` is refused in
// either shape of value: XML keeps it out of an attribute because there is
// nothing an attribute may open, and letting one through would draw a value
// the file did not declare.
func (p *xmlParser) attrValue() (string, error) {
	if p.eof() {
		return "", fmt.Errorf("the value is missing")
	}
	quote := p.src[p.pos]
	if quote != '"' && quote != '\'' {
		// An unquoted value runs to the first space, as XML has it — and not
		// into the slash of a `/>`, which closes the tag rather than being part
		// of what is written on it.
		start := p.pos
		for !p.eof() && !isSpace(p.src[p.pos]) && p.src[p.pos] != '>' && p.src[p.pos] != '<' {
			if p.src[p.pos] == '/' && p.pos+1 < len(p.src) && p.src[p.pos+1] == '>' {
				break
			}
			p.pos++
		}
		return unescape(p.src[start:p.pos]), nil
	}
	p.pos++
	start := p.pos
	for !p.eof() && p.src[p.pos] != quote {
		if p.src[p.pos] == '<' {
			return "", fmt.Errorf("the value holds a `<`, which an attribute may not")
		}
		p.pos++
	}
	if p.eof() {
		return "", fmt.Errorf("the %c that opened the value is never closed", quote)
	}
	v := p.src[start:p.pos]
	p.pos++
	return unescape(v), nil
}
