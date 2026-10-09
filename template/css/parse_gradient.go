package css

// parseGradient reads the arguments of one gradient function into a Gradient.
func parseGradient(name, args string, ctx Units) (*Gradient, bool) {
	switch name {
	case "linear-gradient":
		return parseLinearGradient(args, ctx)
	case "radial-gradient":
		return parseRadialGradient(args, ctx)
	default:
		return parseConicGradient(args, ctx)
	}
}
