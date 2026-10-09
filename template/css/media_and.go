package css

// and conjoins two media lists the way nested @media requires: every pair of
// query alternatives intersects its windows over the viewport. Either side
// without queries passes the other through unchanged.
func (m Media) and(n Media) Media {
	if len(m.Queries) == 0 {
		return n
	}
	if len(n.Queries) == 0 {
		return m
	}
	var out Media
	for _, a := range m.Queries {
		for _, b := range n.Queries {
			out.Queries = append(out.Queries, intersectQuery(a, b))
		}
	}
	return out
}
