package svg

import "fmt"

// tail reads what the text holds after its root element: nothing but
// whitespace, a comment or an instruction may follow, and anything else — a
// second root, a closing tag that belongs to nothing, a word — means the
// drawing is not the single-rooted document it declared itself to be.
func (p *xmlParser) tail() error {
	for {
		p.skipSpace()
		switch {
		case p.eof():
			return nil
		case p.at("<!--"):
			p.through("-->")
		case p.at("<?"):
			p.through("?>")
		default:
			return fmt.Errorf("antui/svg: the drawing carries on past its root with %q, where only whitespace, a comment or an instruction may follow", p.peek(12))
		}
	}
}
