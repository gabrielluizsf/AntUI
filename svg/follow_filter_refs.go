package svg

// followFilterRefs follows what each `url(#name)` in an element's filter list
// names, once the whole drawing is read: the reference was kept as written
// where the list was read, because only a caller holding the drawing can say
// what an id is — the same as a fill, a clip or a mask. What names nothing the
// drawing has is said here, once, and left out of the list the same way a
// function the package cannot read is.
func (img *Image) followFilterRefs(filters []filterOp, warn func(string, ...any)) []filterOp {
	if len(filters) == 0 {
		return filters
	}
	out := make([]filterOp, 0, len(filters))
	for _, f := range filters {
		if f.ref == "" {
			out = append(out, f)
			continue
		}
		def := img.filters[f.ref]
		if def == nil {
			warn("the filter names #%s, which the drawing does not have, so it is left out", f.ref)
			continue
		}
		out = append(out, filterOp{def: def})
	}
	return out
}
