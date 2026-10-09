package css

func blendFilter(a, b Filter, t float64) Filter {
	if a.Drop == nil && b.Drop == nil {
		if a.Kind != b.Kind {
			if t < 0.5 {
				return a
			}
			return b
		}
		return Filter{Kind: a.Kind, Amount: a.Amount + (b.Amount-a.Amount)*t}
	}
	if a.Drop != nil && b.Drop != nil {
		sh := blendShadow(*a.Drop, *b.Drop, t)
		return Filter{Drop: &sh}
	}
	if t < 0.5 {
		return a
	}
	return b
}
