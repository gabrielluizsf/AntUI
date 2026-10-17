package svg

import "fmt"

// checkBinding says whether a prefix may be bound to a namespace at all: the
// `xmlns` prefix is XML's own and may not be bound, and `xml` may only be bound
// to the XML namespace it already stands for.
func checkBinding(prefix, uri string) error {
	if prefix == "xmlns" {
		return fmt.Errorf("antui/svg: the prefix xmlns is XML's own and may not be bound to a namespace")
	}
	if prefix == "xml" && uri != xmlNS {
		return fmt.Errorf("antui/svg: the prefix xml may only be bound to the XML namespace, not to %q", uri)
	}
	return nil
}
