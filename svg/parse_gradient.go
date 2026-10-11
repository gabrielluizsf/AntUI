package svg

// parseGradient reads one `<linearGradient>` or `<radialGradient>`. What it does
// not say is not decided here: an href names another gradient, and only a caller
// holding the whole drawing can say what that one had, so the reference is kept
// and followed by [gradient.resolved].
func parseGradient(e *element, warn func(string, ...any)) *gradient {
	g := gradientDefaults(e.Name == "radialGradient")
	g.href = gradientHref(e.Attr)
	g.stops = gradientStops(e.Kids, warn)
	g.read(e.Attr, warn)
	return g
}
