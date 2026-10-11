package svg

import "strings"

// paintRef is the id inside a `url(#name)` — the thing a fill or a stroke
// points at when it is not naming a colour. The bool says whether there was
// one: `url(#)` with nothing after the hash names nothing, and a paint this
// cannot read is a shape that has to say so rather than come out blank.
//
// A `url(...)` may carry a fallback after it, as in `url(#g) red`: the fallback
// is what the shape is painted with when the reference is not there, so it is
// returned alongside the id for the caller to keep.
func paintRef(s string) (id string, ok bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(strings.ToLower(s), "url(") {
		return "", false
	}
	end := strings.Index(s, ")")
	if end < 0 {
		return "", false
	}
	inside := strings.TrimSpace(s[len("url("):end])
	// Quotes around the reference are allowed and mean nothing.
	inside = strings.Trim(inside, "\"'")
	if !strings.HasPrefix(inside, "#") {
		return "", false
	}
	id = inside[1:]
	return id, id != ""
}
