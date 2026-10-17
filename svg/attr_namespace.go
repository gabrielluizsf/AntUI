package svg

import (
	"fmt"
	"strings"
)

// takeAttrs holds every prefixed attribute to the namespace its prefix is
// bound to. One written against XLink or XML — the two namespaces a drawing's
// attributes use — is kept under the name the drawing reads it by: `xl:href`
// with `xl` bound to XLink is the `xlink:href` a `<use>` looks for. An
// attribute of any other namespace is another document's, and is let go rather
// than read as a name this package happens to know. A prefix nothing declares
// is an error, the same as it is on an element's own name. Two attributes with
// the same local name under prefixes that stand for the same namespace are one
// expanded name written twice, which XML Namespaces keeps out; the `xmlns`
// spellings declare namespaces rather than say anything, and are left as they
// were.
func (ns namespace) takeAttrs(e *element) error {
	var seen map[string]bool
	for key, value := range e.Attr {
		if key == "xmlns" || strings.HasPrefix(key, "xmlns:") {
			continue
		}
		i := strings.IndexByte(key, ':')
		if i <= 0 {
			continue
		}
		prefix, local := key[:i], key[i+1:]
		uri, bound := ns.uriOf(prefix)
		if !bound {
			return fmt.Errorf("antui/svg: <%s> is given %q, where the prefix %q is bound to no namespace", e.Name, key, prefix)
		}
		if seen == nil {
			seen = map[string]bool{}
		}
		if seen[uri+" "+local] {
			return fmt.Errorf("antui/svg: <%s> holds two attributes named %q in the same namespace", e.Name, local)
		}
		seen[uri+" "+local] = true
		delete(e.Attr, key)
		switch uri {
		case xlinkNS:
			e.Attr["xlink:"+local] = value
		case xmlNS:
			e.Attr["xml:"+local] = value
		}
	}
	return nil
}
