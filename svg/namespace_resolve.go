package svg

import (
	"fmt"
	"strings"
)

// uriOf answers what a prefix stands for: the `xml` prefix is bound to the XML
// namespace without being declared, as it has been since namespaces were
// written, and every other prefix has to be declared before it is worn.
func (ns namespace) uriOf(prefix string) (string, bool) {
	if prefix == "xml" {
		return xmlNS, true
	}
	uri, bound := ns[prefix]
	return uri, bound
}

// resolve reads a name through the scope: what it is called once its prefix is
// taken off, and whether it stands outside the SVG namespace, which is the one
// case where what was written is not this package's to draw. A prefix nothing
// declares is an error, the same as it is to an XML Namespaces reader: a
// drawing that wears one is not the drawing it declares, and guessing at what
// it stood for would paint something the file did not say. A default the file
// declares itself is the one the element stands in — an empty one puts it in
// no namespace, which is not the SVG one — while a drawing that declares no
// default at all is read in the SVG namespace, a snippet in its own language.
func (ns namespace) resolve(name string) (string, bool, error) {
	prefix, local := "", name
	if i := strings.IndexByte(name, ':'); i > 0 {
		prefix, local = name[:i], name[i+1:]
	}
	if prefix == "" {
		uri, bound := ns[""]
		return local, bound && uri != svgNS, nil
	}
	uri, bound := ns.uriOf(prefix)
	if !bound {
		return "", false, fmt.Errorf("antui/svg: %q wears the prefix %q, which no namespace binds", name, prefix)
	}
	return local, uri != svgNS, nil
}
