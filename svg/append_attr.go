package svg

import "fmt"

// appendAttr reads one `[…]` attribute condition: the name, and the operator
// and value it is compared with where it asked for one.
func appendAttr(c *compoundSelector, s string, i int) (int, error) {
	i = skipSelSpace(s, i+1)
	name, next := readSelIdent(s, i)
	if name == "" {
		return i, fmt.Errorf("an attribute selector names no attribute")
	}
	i = skipSelSpace(s, next)
	if i < len(s) && s[i] == ']' {
		c.simples = append(c.simples, simpleSelector{kind: kindAttribute, name: name, op: attrExists})
		return i + 1, nil
	}
	op, w := attrOpOf(s[i:])
	if w == 0 {
		return i, fmt.Errorf("an attribute selector has %q where an operator was expected", s[i:])
	}
	val, next, err := readAttrValue(s, i+w)
	if err != nil {
		return i, err
	}
	i = skipSelSpace(s, next)
	if i >= len(s) || s[i] != ']' {
		return i, fmt.Errorf("an attribute selector is never closed with a `]`")
	}
	c.simples = append(c.simples, simpleSelector{kind: kindAttribute, name: name, op: op, val: val})
	return i + 1, nil
}

// attrOpOf reads the operator of an attribute selector and how many bytes it
// took, or a width of none.
func attrOpOf(s string) (attrOp, int) {
	if s == "" {
		return 0, 0
	}
	if s[0] == '=' {
		return attrEquals, 1
	}
	if len(s) > 1 && s[1] == '=' {
		switch s[0] {
		case '~':
			return attrIncludes, 2
		case '|':
			return attrDash, 2
		case '^':
			return attrPrefix, 2
		case '$':
			return attrSuffix, 2
		case '*':
			return attrSubstr, 2
		}
	}
	return 0, 0
}
