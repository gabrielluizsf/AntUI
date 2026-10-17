package svg

import (
	"fmt"
	"strings"
)

// pi reads one `<? … ?>` instruction to its end and leaves the cursor past it.
// An instruction that was never closed is an error rather than the rest of the
// file quietly becoming part of it.
func (p *xmlParser) pi() error {
	p.pos += len("<?")
	i := strings.Index(p.src[p.pos:], "?>")
	if i < 0 {
		return fmt.Errorf("antui/svg: an instruction is never closed")
	}
	p.pos += i + len("?>")
	return nil
}
