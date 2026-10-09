package css

import (
	"fmt"
)

// call evaluates calc(...) and clamp(...). calc is a single expression whose
// additions, subtractions, products and divisions were consumed by expr;
// clamp is three comma-separated arguments with the midpoint clamped by the
// two ends.
func (p *mathParser) call(fn string) (float64, error) {
	if p.t.next() != fn {
		return 0, fmt.Errorf("css: expected %s", fn)
	}
	if p.t.next() != "(" {
		return 0, fmt.Errorf("css: expected (")
	}
	first, err := p.expr()
	if err != nil {
		return 0, err
	}
	if fn == "calc" {
		if p.t.next() != ")" {
			return 0, fmt.Errorf("css: expected )")
		}
		return first, nil
	}
	if p.t.next() != "," {
		return 0, fmt.Errorf("css: expected ,")
	}
	second, err := p.expr()
	if err != nil {
		return 0, err
	}
	if p.t.next() != "," {
		return 0, fmt.Errorf("css: expected ,")
	}
	third, err := p.expr()
	if err != nil {
		return 0, err
	}
	if p.t.next() != ")" {
		return 0, fmt.Errorf("css: expected )")
	}
	lo, hi := first, third
	if lo > hi {
		lo, hi = hi, lo
	}
	if second < lo {
		return lo, nil
	}
	if second > hi {
		return hi, nil
	}
	return second, nil
}
