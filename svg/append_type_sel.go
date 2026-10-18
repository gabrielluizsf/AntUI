package svg

import "fmt"

// appendTypeSel reads a type selector, with a namespace prefix where one is
// written. The prefix is resolved against the namespaces in scope where the
// `<style>` stands: the elements a drawing keeps are all in the SVG namespace,
// so a type naming any other is read and never matches, and one wearing a
// prefix nothing binds is not a sheet this drawing can read.
func appendTypeSel(c *compoundSelector, s string, i int, ns namespace) (int, error) {
	name, next := readSelIdent(s, i)
	if next < len(s) && s[next] == '|' && (next+1 >= len(s) || s[next+1] != '|') {
		return appendLocalType(c, s, next+1, name, ns)
	}
	c.simples = append(c.simples, simpleSelector{kind: kindType, name: name})
	return next, nil
}

// appendLocalType reads the type after a namespace prefix: a name, or `*` for
// every name in that namespace. The prefix is the one written before the `|`.
func appendLocalType(c *compoundSelector, s string, i int, prefix string, ns namespace) (int, error) {
	uri, err := ns.selNS(prefix)
	if err != nil {
		return i, err
	}
	stands := standsWithSVG(uri)
	if i < len(s) && s[i] == '*' {
		if stands {
			c.simples = append(c.simples, simpleSelector{kind: kindUniversal})
		} else {
			c.simples = append(c.simples, simpleSelector{kind: kindUnmatchable})
		}
		return i + 1, nil
	}
	name, next := readSelIdent(s, i)
	if name == "" {
		return i, fmt.Errorf("a namespace prefix names no type")
	}
	if !stands {
		c.simples = append(c.simples, simpleSelector{kind: kindUnmatchable})
		return next, nil
	}
	c.simples = append(c.simples, simpleSelector{kind: kindType, name: name})
	return next, nil
}
