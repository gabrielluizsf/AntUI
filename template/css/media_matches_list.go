package css

// matches reports whether any query alternative holds at this viewport. An
// empty query list (a rule with no @media at all) always matches.
func (m Media) matches(vp Viewport) bool {
	if len(m.Queries) == 0 {
		return true
	}
	for _, q := range m.Queries {
		if q.matches(vp) {
			return true
		}
	}
	return false
}
