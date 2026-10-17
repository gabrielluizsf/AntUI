package svg

import (
	"fmt"
	"strings"
)

// pi reads one `<? … ?>` instruction to its end and leaves the cursor past it.
// An instruction that was never closed is an error rather than the rest of the
// file quietly becoming part of it, and the target `xml` — however it is
// spelled — is an error too: XML reserves it for the declaration that may only
// open a document, and an instruction wearing it says nothing a reader can
// take.
func (p *xmlParser) pi() error {
	p.pos += len("<?")
	if strings.EqualFold(p.readName(), "xml") {
		return fmt.Errorf("antui/svg: the target xml is reserved for the declaration, which may only open a document")
	}
	i := strings.Index(p.src[p.pos:], "?>")
	if i < 0 {
		return fmt.Errorf("antui/svg: an instruction is never closed")
	}
	p.pos += i + len("?>")
	return nil
}
