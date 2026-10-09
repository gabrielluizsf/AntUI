package css

// setProp copies one canonical property's value from src to dst and marks it
// in dst's Set map. Geometric properties land in setPropBox and the paint
// properties in setPropFX, each case table staying under the file ceiling.
func setProp(dst *Style, src Style, prop string) {
	switch prop {
	case "background-color", "color", "opacity", "transform", "transform-origin",
		"box-shadow", "text-shadow", "filter", "backdrop-filter":
		setPropFX(dst, src, prop)
	default:
		setPropBox(dst, src, prop)
	}
	if dst.Set != nil {
		dst.Set[prop] = true
	}
}
