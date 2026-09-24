package css

import (
	"fmt"
	"strconv"
	"strings"
)

// EvalMath is the length-formula engine behind calc(), min(), max() and
// clamp(): it turns "calc(100% - 20px)", "min(480px, 50vw)" and friends into
// reference pixels under a measurement context. Every term collapses to
// pixels — percentages against the containing width, viewport and font units
// against the relevant measure — and the arithmetic runs in that space, which
// is exactly how a browser expands a length mix. ok is false when the formula
// does not balance or a term does not parse, and the caller warns or drops.
func EvalMath(raw string, ctx Units) (px float64, ok bool) {
	t := &mathTokens{s: stripComments(raw)}
	p := &mathParser{t: t, ctx: ctx}
	v, err := p.expr()
	if err != nil || !t.empty() {
		return 0, false
	}
	return v, true
}

// mathLook reports whether a raw value opens a math formula the engine can
// evaluate (calc/min/max/clamp). Anything else with a "(" in the way is a
// function this engine does not understand, not something to attempt.
func mathLook(raw string) bool {
	t := strings.ToLower(strings.TrimSpace(stripComments(raw)))
	for _, fn := range []string{"calc(", "min(", "max(", "clamp("} {
		if strings.HasPrefix(t, fn) {
			return true
		}
	}
	return false
}

// mathParser is a small recursive-descent evaluator over + - * / ( ) ,
// numbers-with-units and the four math functions.
type mathParser struct {
	t   *mathTokens
	ctx Units
}

// expr parses a sum of terms.
func (p *mathParser) expr() (float64, error) {
	v, err := p.term()
	if err != nil {
		return 0, err
	}
	for {
		switch p.t.peek() {
		case "+", "-":
			op := p.t.next()
			w, err := p.term()
			if err != nil {
				return 0, err
			}
			if op == "+" {
				v += w
			} else {
				v -= w
			}
		default:
			return v, nil
		}
	}
}

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

// number reads a numeric literal with an optional unit, converted to
// reference pixels under the parser's context. A plain number is pixels.
func (p *mathParser) number() (float64, error) {
	tok := p.t.next()
	if tok == "" {
		return 0, fmt.Errorf("css: expected a number")
	}
	// The longest matching unit wins — rem before em, vmin before vw, etc.
	u := unitPx
	for _, suf := range unitSuffixes {
		if strings.HasSuffix(tok, suf.name) {
			u = suf.u
			tok = strings.TrimSuffix(tok, suf.name)
			break
		}
	}
	if strings.HasSuffix(tok, "%") {
		u = unitPct
		tok = strings.TrimSuffix(tok, "%")
	}
	v, err := strconv.ParseFloat(tok, 64)
	if err != nil {
		return 0, fmt.Errorf("css: %q is not a number", tok)
	}
	l := Length{u: u, value: v}
	return l.ref(p.ctx), nil
}

// mathTokens splits a formula into operators, parentheses, commas and
// numbers-with-units. "calc(" and "min(" stay whole so the parser matches the
// function name before its "(".
type mathTokens struct {
	s    string
	toks []string
	i    int
}

func (t *mathTokens) split() []string {
	if t.toks != nil {
		return t.toks
	}
	s := t.s
	var out []string
	for i := 0; i < len(s); {
		c := s[i]
		switch {
		case isSpace(c):
			i++
		case strings.IndexByte("+-*/(),", c) >= 0:
			// A function name ("calc", "min"… ) is the letters just before
			// its "(", handed to the parser as one token.
			if c == '(' {
				name := ""
				j := i - 1
				for j >= 0 && isIdentByte(s[j]) {
					name = string(s[j]) + name
					j--
				}
				switch name {
				case "calc", "min", "max", "clamp":
					out = append(out, name)
					out = append(out, "(")
					i++
					continue
				}
			}
			out = append(out, string(c))
			i++
		case c == '.' || (c >= '0' && c <= '9'):
			j := i
			for j < len(s) && (s[j] == '.' || (s[j] >= '0' && s[j] <= '9')) {
				j++
			}
			for j < len(s) && isUnitByte(s[j]) {
				j++
			}
			out = append(out, s[i:j])
			i = j
		default:
			// A stray glyph: skip it rather than fail the whole formula.
			i++
		}
	}
	t.toks = out
	return out
}

func (t *mathTokens) peek() string {
	ts := t.split()
	if t.i >= len(ts) {
		return ""
	}
	return ts[t.i]
}

func (t *mathTokens) next() string {
	ts := t.split()
	if t.i >= len(ts) {
		return ""
	}
	v := ts[t.i]
	t.i++
	return v
}

func (t *mathTokens) empty() bool {
	return t.peek() == ""
}

// isUnitByte reports whether a byte can continue a numeric literal's unit.
func isUnitByte(c byte) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c == '%':
		return true
	}
	return false
}
