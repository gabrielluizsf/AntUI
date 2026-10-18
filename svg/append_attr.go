package svg

import "fmt"

// appendAttr reads one `[…]` attribute condition: the name, the operator and
// value it is compared with where it asked for one, and the `i` or `s` flag
// that may follow. The name's namespace is resolved as a type selector's is,
// so a sheet can tell the `href` of no namespace from XLink's.
func appendAttr(c *compoundSelector, s string, i int, ns namespace) (int, error) {
	i = skipSelSpace(s, i+1)
	name, prefix, named, next := readAttrName(s, i)
	if name == "" {
		return i, fmt.Errorf("an attribute selector names no attribute")
	}
	uri, any, err := attrNS(prefix, named, ns)
	if err != nil {
		return i, err
	}
	cond := simpleSelector{kind: kindAttribute, name: name, uri: uri, any: any}
	i = skipSelSpace(s, next)
	if i < len(s) && s[i] == ']' {
		cond.op = attrExists
		c.simples = append(c.simples, cond)
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
	caseless, i, err := attrCaseFlag(s, i)
	if err != nil {
		return i, err
	}
	if i >= len(s) || s[i] != ']' {
		return i, fmt.Errorf("an attribute selector is never closed with a `]`")
	}
	cond.op, cond.val, cond.caseless = op, val, caseless
	c.simples = append(c.simples, cond)
	return i + 1, nil
}
