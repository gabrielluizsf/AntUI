package svg

import "strings"

// namespace is what the prefixes in scope at one place in a drawing stand for:
// the empty prefix is the default one, which an element without a prefix of
// its own is read through. The SVG namespace is the only one this package
// draws; an element standing anywhere else is a drawing within the drawing,
// and what is written inside it is not this package's to paint.
type namespace map[string]string

// svgNS is the namespace an SVG element stands in, written out full.
const svgNS = "http://www.w3.org/2000/svg"

// child answers the scope an element and everything written inside it is read
// through: its own, unless the element binds a prefix or the default, which is
// only copied into a map of its own when there is something to copy — a
// drawing binds its namespaces once, at the root, and every element below
// shares what came before it.
func (ns namespace) child(e *element) namespace {
	var inner namespace
	for key, value := range e.Attr {
		if key != "xmlns" && !strings.HasPrefix(key, "xmlns:") {
			continue
		}
		if inner == nil {
			inner = namespace{}
			for prefix, uri := range ns {
				inner[prefix] = uri
			}
		}
		prefix := ""
		if key != "xmlns" {
			prefix = key[len("xmlns:"):]
		}
		inner[prefix] = value
	}
	if inner == nil {
		return ns
	}
	return inner
}

// resolve reads a name through the scope: what it is called once its prefix is
// taken off, and whether it stands outside the SVG namespace, which is the one
// case where what was written is not this package's to draw. A prefix nothing
// declares is let through by its local name, because that is how the drawings
// that reach here spell one; an XML Namespaces reader would call that an
// error, and the day the drawings change their minds this is the line to fix.
func (ns namespace) resolve(name string) (string, bool) {
	prefix, local := "", name
	if i := strings.IndexByte(name, ':'); i > 0 {
		prefix, local = name[:i], name[i+1:]
	}
	uri, bound := ns[prefix]
	if !bound {
		return local, false
	}
	return local, uri != svgNS
}
