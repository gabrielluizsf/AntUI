package svg

import "strings"

// through moves past the next occurrence of what, and past the end of the
// document if it is not there, which is what a comment or a CDATA section that
// was never finished wants: the rest of the file is not worth reading.
func (p *xmlParser) through(what string) {
	if i := strings.Index(p.src[p.pos:], what); i >= 0 {
		p.pos += i + len(what)
		return
	}
	p.pos = len(p.src)
}
