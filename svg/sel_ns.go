package svg

import "fmt"

// selNS answers the namespace a selector's prefix stands for. `*` is every
// namespace at once; every other prefix has to be bound before it is worn, the
// same as in a document read by an XML Namespaces reader, so a sheet that
// wears one nothing declares is not a sheet this drawing can read.
func (ns namespace) selNS(prefix string) (string, error) {
	if prefix == "*" {
		return "*", nil
	}
	uri, bound := ns.uriOf(prefix)
	if !bound {
		return "", fmt.Errorf("the namespace prefix %q is one nothing binds", prefix)
	}
	return uri, nil
}

// standsWithSVG answers whether a selector naming one namespace may stand on
// the elements a drawing keeps. Every element of one is in the SVG namespace —
// an element standing anywhere else is a drawing within the drawing and is
// left out of the tree — so a selector naming another is read all the same
// and simply never matches.
func standsWithSVG(uri string) bool { return uri == "*" || uri == svgNS }
