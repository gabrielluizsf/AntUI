package svg

import "fmt"

// doctypeStep reads one piece of the doctype — a comment, an instruction, a
// quoted part, or a bracket that opens or closes the internal subset — and
// answers whether the doctype ended, what the subset spans after it, and the
// first of those that could not be read.
func (p *xmlParser) doctypeStep(from int) (bool, int, error) {
	switch {
	case p.at("<!--"):
		return false, from, p.comment()
	case p.at("<?"):
		return false, from, p.pi()
	case p.at("["):
		p.pos++
		return false, p.pos, nil
	case p.at("]"):
		if from < 0 {
			return false, from, fmt.Errorf("antui/svg: a `]` in a doctype comes with no subset to close")
		}
		if err := p.parseSubset(from, p.pos); err != nil {
			return false, from, err
		}
		p.pos++
		return false, -1, nil
	case p.src[p.pos] == '"' || p.src[p.pos] == '\'':
		return false, from, p.inQuotes()
	case p.at(">") && from < 0:
		p.pos++
		return true, from, nil
	}
	p.pos++
	return false, from, nil
}
