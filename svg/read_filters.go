package svg

import "strings"

// readFilters reads the list a `filter` attribute holds: `none` is no filter
// at all and says so without complaint, and every function is read on its own,
// so a list this package cannot wholly read still does the part of it that it
// can — one function left out is a warning, not a whole filter dropped.
func readFilters(raw string, warn func(string, ...any)) []filterOp {
	v := strings.TrimSpace(raw)
	if v == "" {
		warn("the filter %q is not one this package can read, so it is left out", raw)
		return nil
	}
	if strings.EqualFold(v, "none") {
		return nil
	}
	var out []filterOp
	for _, tok := range filterTokens(v) {
		if isFilterRef(tok) {
			// A reference to a `<filter>` element elsewhere in the drawing is
			// the other way a file asks for one of these. The id is kept as
			// written — only a caller holding the whole drawing can say what
			// an id in the file is — and the rest of the list goes on. What
			// names nothing the drawing has is said once the drawing is all
			// read; see [Image.followFilterRefs].
			if id, ok := paintRef(tok); ok {
				out = append(out, filterOp{ref: id})
				continue
			}
			warn("the filter reference %q is not a reference to a <filter> this package can follow, so it is left out", tok)
			continue
		}
		if op, ok := readFilter(tok); ok {
			out = append(out, op)
			continue
		}
		warn("the filter %q is not one this package can read, so it is left out", tok)
	}
	return out
}
