package css

// BlendStyles interpolates a named set of properties from a to b and returns
// a style holding the result, with b's values for everything else.
func BlendStyles(a, b Style, t float64, props map[string]bool, ctx Units) Style {
	out := CopyStyle(b)
	for prop := range props {
		blendProp(&out, a, prop, t, ctx)
	}
	return out
}
