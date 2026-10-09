package css

// setPropFX copies one paint property's value from src to dst: the colours,
// the transform and the shadow and filter lists.
func setPropFX(dst *Style, src Style, prop string) {
	switch prop {
	case "background-color":
		dst.Background = src.Background
	case "color":
		dst.Color = src.Color
	case "opacity":
		dst.Opacity = src.Opacity
	case "transform":
		dst.Transform = src.Transform
	case "transform-origin":
		dst.TransformOrigin = src.TransformOrigin
	case "box-shadow":
		dst.BoxShadow = src.BoxShadow
	case "text-shadow":
		dst.TextShadow = src.TextShadow
	case "filter":
		dst.Filters = src.Filters
	case "backdrop-filter":
		dst.BackdropFilters = src.BackdropFilters
	}
}
