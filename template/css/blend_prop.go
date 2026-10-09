package css

// blendProp interpolates one canonical property from a to b, writing the
// result into b (and its Set map): dst starts as the target value, which is
// exactly how a blend composes. Discrete properties swap over at the half
// way point.
func blendProp(dst *Style, from Style, prop string, t float64, ctx Units) {
	switch prop {
	case "background-color", "color", "opacity", "transform", "transform-origin",
		"box-shadow", "text-shadow", "filter", "backdrop-filter":
		blendPropFX(dst, from, prop, t, ctx)
	default:
		blendPropBox(dst, from, prop, t, ctx)
	}
	if dst.Set != nil {
		dst.Set[prop] = true
	}
}
