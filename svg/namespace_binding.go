package svg

import "fmt"

// checkBinding says whether a prefix may be bound to a namespace at all: the
// `xmlns` prefix is XML's own and may not be bound, `xml` may only be bound to
// the XML namespace it already stands for, and any other prefix has to be
// bound to something — a prefixed declaration with no namespace is an error,
// while only the default one may be emptied to say an element stands in none.
func checkBinding(prefix, uri string) error {
	if prefix == "xmlns" {
		return fmt.Errorf("antui/svg: the prefix xmlns is XML's own and may not be bound to a namespace")
	}
	if prefix == "xml" && uri != xmlNS {
		return fmt.Errorf("antui/svg: the prefix xml may only be bound to the XML namespace, not to %q", uri)
	}
	if prefix != "" && uri == "" {
		return fmt.Errorf("antui/svg: the prefix %q is declared to no namespace, where a prefix needs one", prefix)
	}
	return nil
}
