package svg

import "strings"

// overflowIn is whether an element lets what it draws run past the room it
// asked for — the room of a `<marker>` or the tile of a `<pattern>`, which are
// the two places the spec says `overflow` stands. Only `visible` says so:
// anything else, including nothing at all, says the room cuts what runs past
// its edge, which is what the spec's own stylesheet writes on both of them.
// It is read the way clip-rule is — an attribute first, then a declaration in
// the `style` — and, like clip-rule, read without a warn: a value that is not
// `visible` is simply the cut.
func overflowIn(e *element) bool {
	if e.hasAttr("overflow") {
		return strings.EqualFold(strings.TrimSpace(e.attr("overflow")), "visible")
	}
	if s := e.attr("style"); s != "" {
		if raw, ok := parseDeclarations(s)["overflow"]; ok {
			return strings.EqualFold(strings.TrimSpace(raw), "visible")
		}
	}
	return false
}
