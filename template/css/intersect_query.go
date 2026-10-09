package css

// intersectQuery narrows two windows over the viewport to their overlap, edge
// by edge. Negation is left on the outer window's side only when it is safe to
// guess, which in practice means: keep whichever side actually constrains the
// viewport, and a pair where the negations disagree degrades to "always
// matches" (reported, not fatal).
func intersectQuery(a, b MediaQuery) MediaQuery {
	if a.Negated != b.Negated {
		// A negated window AND a normal one is rare; fall back to the normal
		// side so the rule is not silently dropped.
		if !a.Negated {
			return a
		}
		return b
	}
	q := MediaQuery{Negated: a.Negated}
	q.MinWidth, q.HasMin = intersectBound(a.HasMin, a.MinWidth, b.HasMin, b.MinWidth, true)
	q.MaxWidth, q.HasMax = intersectBound(a.HasMax, a.MaxWidth, b.HasMax, b.MaxWidth, false)
	q.MinHeight, q.HasMinHeight = intersectBound(a.HasMinHeight, a.MinHeight, b.HasMinHeight, b.MinHeight, true)
	q.MaxHeight, q.HasMaxHeight = intersectBound(a.HasMaxHeight, a.MaxHeight, b.HasMaxHeight, b.MaxHeight, false)
	q.MinDpi, q.HasMinDpi = intersectBound(a.HasMinDpi, a.MinDpi, b.HasMinDpi, b.MinDpi, true)
	q.MaxDpi, q.HasMaxDpi = intersectBound(a.HasMaxDpi, a.MaxDpi, b.HasMaxDpi, b.MaxDpi, false)
	// Orientation and the color scheme are the two features that are not
	// ranges. Where the nested queries agree, the answer carries over; where
	// they disagree there is no window on both sides of it, and the rule
	// belongs to nothing.
	q.HasOrientation = a.HasOrientation || b.HasOrientation
	if a.HasOrientation {
		q.Portrait = a.Portrait
	}
	if b.HasOrientation && a.HasOrientation && a.Portrait != b.Portrait {
		q.never = true
	}
	if b.HasOrientation && !a.HasOrientation {
		q.Portrait = b.Portrait
	}
	q.HasScheme = a.HasScheme || b.HasScheme
	if a.HasScheme {
		q.Dark = a.Dark
	}
	if b.HasScheme && a.HasScheme && a.Dark != b.Dark {
		q.never = true
	}
	if b.HasScheme && !a.HasScheme {
		q.Dark = b.Dark
	}
	return q
}
