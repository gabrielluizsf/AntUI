package svg

import "fmt"

// spendExpansion counts n more bytes read out of the doctype's names against
// [maxEntityBytes], and is the error that refuses the drawing once they are
// spent.
func (p *xmlParser) spendExpansion(n int) error {
	p.expanded += n
	if p.expanded > maxEntityBytes {
		return fmt.Errorf("the doctype's names read into more than %d bytes of text", maxEntityBytes)
	}
	return nil
}
