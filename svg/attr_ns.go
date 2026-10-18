package svg

// attrNS reads the namespace an attribute selector asks for: `*|` stands for
// every namespace, a bare name or `|` for the one an unprefixed attribute is
// in — no namespace at all — and any other prefix has to be bound before it is
// worn, as in a document read by an XML Namespaces reader.
func attrNS(prefix string, named bool, ns namespace) (string, bool, error) {
	if named && prefix == "*" {
		return "", true, nil
	}
	if !named || prefix == "" {
		return "", false, nil
	}
	uri, err := ns.selNS(prefix)
	if err != nil {
		return "", false, err
	}
	return uri, false, nil
}
