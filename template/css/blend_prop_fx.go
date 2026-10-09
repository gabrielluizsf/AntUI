package css

import "github.com/gabrielluizsf/antui/canvas"

// blendPropFX interpolates the paint properties: the colours mix, the
// transform, shadows and filters blend through their own list rules, and
// opacity lerps.
func blendPropFX(dst *Style, from Style, prop string, t float64, ctx Units) {
	switch prop {
	case "background-color":
		dst.Background = canvas.Mix(from.Background, dst.Background, float32(t))
	case "color":
		dst.Color = canvas.Mix(from.Color, dst.Color, float32(t))
	case "opacity":
		dst.Opacity = from.Opacity + (dst.Opacity-from.Opacity)*t
	case "transform":
		dst.Transform = blendTransformList(from.Transform, dst.Transform, t)
	case "transform-origin":
		dst.TransformOrigin = [2]Length{
			blendLength(from.TransformOrigin[0], dst.TransformOrigin[0], t, ctx),
			blendLength(from.TransformOrigin[1], dst.TransformOrigin[1], t, ctx),
		}
	case "box-shadow":
		dst.BoxShadow = blendShadowList(from.BoxShadow, dst.BoxShadow, t)
	case "text-shadow":
		dst.TextShadow = blendShadowList(from.TextShadow, dst.TextShadow, t)
	case "filter":
		dst.Filters = blendFilterList(from.Filters, dst.Filters, t)
	case "backdrop-filter":
		dst.BackdropFilters = blendFilterList(from.BackdropFilters, dst.BackdropFilters, t)
	}
}
