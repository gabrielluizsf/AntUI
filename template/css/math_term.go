package css

import (
	"fmt"
)

// term parses a product of factors.
func (p *mathParser) term() (float64, error) {
	v, err := p.factor()
	if err != nil {
		return 0, err
	}
	for {
		switch p.t.peek() {
		case "*", "/":
			op := p.t.next()
			w, err := p.factor()
			if err != nil {
				return 0, err
			}
			if op == "*" {
				v *= w
			} else {
				if w == 0 {
					return 0, fmt.Errorf("css: division by zero")
				}
				v /= w
			}
		default:
			return v, nil
		}
	}
}
