package svg

import "strings"

// readBase is the base an element's relative addresses resolve against, which
// `xml:base` sets: the address it names is itself resolved against the base
// already in hand, so a chain of them reads the way XML Base says, and an
// element that says nothing leaves the inherited base where it was. Anything
// the attribute names that is not blank is taken at its word.
func readBase(raw, inherited string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return inherited
	}
	return resolveHref(inherited, raw)
}
