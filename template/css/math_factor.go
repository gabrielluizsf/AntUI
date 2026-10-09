package css

import (
	"fmt"
)

// factor parses a number, a signed factor, a parenthesised expression or a
// math function. Unary minus is a factor of its own, so calc(-10px) and
// calc(-10px + 50%) resolve the way CSS writes them; the expr level consumes
// the binary operator, leaving factor to see only a leading sign.
func (p *mathParser) factor() (float64, error) {
	tok := p.t.peek()
	switch {
	case tok == "-", tok == "+":
		sign := 1.0
		if tok == "-" {
			sign = -1
		}
		p.t.next()
		v, err := p.factor()
		if err != nil {
			return 0, err
		}
		return sign * v, nil
	case tok == "(":
		p.t.next()
		v, err := p.expr()
		if err != nil {
			return 0, err
		}
		if p.t.next() != ")" {
			return 0, fmt.Errorf("css: unbalanced parentheses")
		}
		return v, nil
	case tok == "calc":
		return p.call("calc")
	case tok == "min":
		return p.minmax(false)
	case tok == "max":
		return p.minmax(true)
	case tok == "clamp":
		return p.call("clamp")
	}
	return p.number()
}
