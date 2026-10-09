package css

// matches tests one alternative against the viewport: everything it declares
// must hold, and a negated query flips that answer as a whole. A query that
// contradicted its own nesting matches nothing at all, and one whose sensors
// all withheld their answer matches everything — a condition the canvas
// cannot judge is neither true nor false for it, so the rule stands.
func (q MediaQuery) matches(vp Viewport) bool {
	if q.never {
		return false
	}
	if !q.constrained() {
		return true
	}
	in := true
	answered := false
	if q.HasMin && vp.Width < q.MinWidth {
		in = false
	}
	if q.HasMax && vp.Width > q.MaxWidth {
		in = false
	}
	if q.HasMinHeight && vp.Height < q.MinHeight {
		in = false
	}
	if q.HasMaxHeight && vp.Height > q.MaxHeight {
		in = false
	}
	if q.HasMin || q.HasMax || q.HasMinHeight || q.HasMaxHeight {
		answered = true
	}
	if q.HasOrientation {
		answered = true
		if q.Portrait != (vp.Height >= vp.Width) {
			in = false
		}
	}
	if dpi, ok := vp.resolutionDpi(); ok && (q.HasMinDpi || q.HasMaxDpi) {
		answered = true
		if q.HasMinDpi && dpi < q.MinDpi {
			in = false
		}
		if q.HasMaxDpi && dpi > q.MaxDpi {
			in = false
		}
	}
	if q.HasScheme && vp.Scheme != SchemeUnknown {
		answered = true
		if q.Dark != (vp.Scheme == SchemeDark) {
			in = false
		}
	}
	if !answered {
		return true
	}
	if q.Negated {
		return !in
	}
	return in
}
