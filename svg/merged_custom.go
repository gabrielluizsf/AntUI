package svg

// mergedCustom is what var() sees on an element: the custom properties written
// above it with the ones it declares itself on top, which is the inheritance a
// custom property is given. The inherited map is copied rather than written
// through, because the style it came from is the one every sibling reads.
func mergedCustom(inherited, declared map[string]string) map[string]string {
	out := make(map[string]string, len(inherited)+len(declared))
	for name, raw := range inherited {
		out[name] = raw
	}
	for name, raw := range declared {
		out[name] = raw
	}
	return out
}
