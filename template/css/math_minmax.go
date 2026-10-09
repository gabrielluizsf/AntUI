package css

import (
	"fmt"
)

// minmax evaluates min(...) or max(...): a comma-separated list, the smallest
// or the largest wins.
func (p *mathParser) minmax(isMax bool) (float64, error) {
	fn := "min"
	if isMax {
		fn = "max"
	}
	if p.t.next() != fn {
		return 0, fmt.Errorf("css: expected %s", fn)
	}
	if p.t.next() != "(" {
		return 0, fmt.Errorf("css: expected (")
	}
	v, err := p.expr()
	if err != nil {
		return 0, err
	}
	for p.t.peek() == "," {
		p.t.next()
		w, err := p.expr()
		if err != nil {
			return 0, err
		}
		if isMax {
			v = max(v, w)
		} else {
			v = min(v, w)
		}
	}
	if p.t.next() != ")" {
		return 0, fmt.Errorf("css: expected )")
	}
	return v, nil
}
