package svg

// attrUnder answers the value of the attribute a selector names, in the
// namespace it asked for. An attribute of no namespace is written plain; the
// two a drawing reads — the ones XLink and XML bind — carry their own prefix;
// and `*|` stands for either.
func attrUnder(e *element, s simpleSelector) (string, bool) {
	if s.any {
		for _, key := range []string{s.name, "xlink:" + s.name, "xml:" + s.name} {
			if v, ok := e.Attr[key]; ok {
				return v, true
			}
		}
		return "", false
	}
	switch s.uri {
	case xlinkNS:
		v, ok := e.Attr["xlink:"+s.name]
		return v, ok
	case xmlNS:
		v, ok := e.Attr["xml:"+s.name]
		return v, ok
	case "":
		v, ok := e.Attr[s.name]
		return v, ok
	}
	// Every other namespace: none a drawing keeps, so none to find.
	return "", false
}
