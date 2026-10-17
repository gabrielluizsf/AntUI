package svg

import (
	"fmt"
	"strings"
)

// attrSpace reads the whitespace XML normalizes inside every attribute value:
// a tab or a line break written literally in the value is one space. A
// character reference holds the character itself and is left to be read as it
// says, which is why this runs before the value is unescaped.
var attrSpace = strings.NewReplacer("\t", " ", "\n", " ")

// attrValue reads the value of one attribute, unwrapping its quotes and turning
// what is inside them into the characters it stands for. A value is taken in
// quotes only, as XML has it — one spelled without them is not a value a reader
// can take, since where it ends is anyone's guess. A `<` is refused inside the
// quotes too: XML keeps it out of an attribute because there is nothing an
// attribute may open, and letting one through would draw a value the file did
// not declare.
func (p *xmlParser) attrValue() (string, error) {
	if p.eof() {
		return "", fmt.Errorf("the value is missing")
	}
	quote := p.src[p.pos]
	if quote != '"' && quote != '\'' {
		return "", fmt.Errorf("the value is not in quotes, which every attribute value is")
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
	return p.unescape(attrSpace.Replace(v))
}
