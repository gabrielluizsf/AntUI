package svg

import "strings"

// stripBOM takes the byte-order mark Windows editors leave at the front of a
// file. It is not part of the drawing and not a name, so without this a drawing
// that opens with one is refused before its first tag is even looked for.
func (p *xmlParser) stripBOM() {
	p.src = strings.TrimPrefix(p.src, "\ufeff")
}
