package svg

import (
	"fmt"
	"strings"
)

// comment reads a `<!-- … -->` to its end and leaves the cursor past it. A
// comment that was never closed is an error rather than the rest of the file
// quietly becoming invisible, and so is one holding `--` inside it or ending in
// a dash: XML keeps those out because there is no telling where such a
// comment was meant to end.
func (p *xmlParser) comment() error {
	p.pos += len("<!--")
	i := strings.Index(p.src[p.pos:], "-->")
	if i < 0 {
		return fmt.Errorf("antui/svg: a comment is never closed")
	}
	body := p.src[p.pos : p.pos+i]
	if strings.Contains(body, "--") || strings.HasSuffix(body, "-") {
		return fmt.Errorf("antui/svg: a comment holds a run of dashes, which XML keeps out of one")
	}
	p.pos += i + len("-->")
	return nil
}
