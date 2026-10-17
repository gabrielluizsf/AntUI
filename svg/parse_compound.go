package svg

import "fmt"

// parseCompound reads the simple selectors written with nothing between them,
// stopping at the first byte that begins a combinator, a comma or the end.
func parseCompound(s string, i int) (compoundSelector, int, error) {
	var c compoundSelector
	for i < len(s) {
		var err error
		switch ch := s[i]; {
		case ch == '*':
			c.simples = append(c.simples, simpleSelector{kind: kindUniversal})
			i++
		case ch == '.':
			i, err = appendNamed(&c, kindClass, s, i+1)
		case ch == '#':
			i, err = appendNamed(&c, kindID, s, i+1)
		case ch == '[':
			i, err = appendAttr(&c, s, i)
		case ch == ':':
			i, err = appendPseudo(&c, s, i+1)
		case isNameByte(ch, true):
			i, err = appendNamed(&c, kindType, s, i)
		default:
			if len(c.simples) == 0 {
				return c, i, fmt.Errorf("a selector starts with %q, which is not one", s[i:])
			}
			return c, i, nil
		}
		if err != nil {
			return c, i, err
		}
	}
	if len(c.simples) == 0 {
		return c, i, fmt.Errorf("a selector is empty")
	}
	return c, i, nil
}
