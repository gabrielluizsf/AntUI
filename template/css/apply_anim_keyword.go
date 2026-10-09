package css

// applyAnimKeyword sets the field a direction or fill keyword names and
// reports whether it recognised one.
func applyAnimKeyword(a *Animation, t string) bool {
	switch t {
	case "infinite":
		a.Infinite = true
	case "reverse":
		a.Direction = AnimReverse
	case "alternate":
		a.Direction = AnimAlternate
	case "alternate-reverse":
		a.Direction = AnimAlternateReverse
	case "backwards":
		a.Fill = FillBackwards
	case "forwards":
		a.Fill = FillForwards
	case "both":
		a.Fill = FillBoth
	case "none":
		// fill-mode: none, the default.
	default:
		return false
	}
	return true
}
