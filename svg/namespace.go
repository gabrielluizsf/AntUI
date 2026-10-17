package svg

import "strings"

// namespace is what the prefixes in scope at one place in a drawing stand for:
// the empty prefix is the default one, which an element without a prefix of
// its own is read through. The SVG namespace is the only one this package
// draws; an element standing anywhere else is a drawing within the drawing,
// and what is written inside it is not this package's to paint.
type namespace map[string]string

// The namespaces a drawing may lean on without declaring them, and the one the
// SVG elements themselves stand in.
const (
	svgNS   = "http://www.w3.org/2000/svg"
	xmlNS   = "http://www.w3.org/XML/1998/namespace"
	xlinkNS = "http://www.w3.org/1999/xlink"
)

// child answers the scope an element and everything written inside it is read
// through: its own, unless the element binds a prefix or the default, which is
// only copied into a map of its own when there is something to copy — a
// drawing binds its namespaces once, at the root, and every element below
// shares what came before it. The two prefixes XML keeps for itself are not a
// file's to bind to anything else.
func (ns namespace) child(e *element) (namespace, error) {
	var inner namespace
	for key, value := range e.Attr {
		if key != "xmlns" && !strings.HasPrefix(key, "xmlns:") {
			continue
		}
		prefix := ""
		if key != "xmlns" {
			prefix = key[len("xmlns:"):]
		}
		if err := checkBinding(prefix, value); err != nil {
			return nil, err
		}
		if inner == nil {
			inner = namespace{}
			for p, uri := range ns {
				inner[p] = uri
			}
		}
		inner[prefix] = value
	}
	if inner == nil {
		return ns, nil
	}
	return inner, nil
}
