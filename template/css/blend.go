package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// animatableProps is every canonical property the engine interpolates
// between two styles. It is the cylinder a transition-property of "all" and
// a keyframe both draw from; a property absent here snaps instead. Border
// widths and the discrete layout keywords (display, position, box-sizing,
// visibility) deliberately stay away.
var animatableProps = []string{
	"background-color", "color", "opacity",
	"width", "height",
	"min-width", "max-width", "min-height", "max-height",
	"margin-top", "margin-right", "margin-bottom", "margin-left",
	"padding-top", "padding-right", "padding-bottom", "padding-left",
	"top", "right", "bottom", "left",
	"font-size", "letter-spacing", "word-spacing", "line-height",
	"transform", "transform-origin",
	"box-shadow", "text-shadow", "filter", "backdrop-filter",
}

// CopyStyle returns a style that owns its Set map, so a caller can blend or
// freeze values into it without touching the cached cascade's Set. The
// Custom and inherit maps are read-only after the cascade and stay shared.
func CopyStyle(s Style) Style {
	out := make(map[string]bool, len(s.Set))
	for k := range s.Set {
		out[k] = true
	}
	s.Set = out
	return s
}

// BlendStyles interpolates a named set of properties from a to b and returns
// a style holding the result, with b's values for everything else.
func BlendStyles(a, b Style, t float64, props map[string]bool, ctx Units) Style {
	out := CopyStyle(b)
	for prop := range props {
		blendProp(&out, a, prop, t, ctx)
	}
	return out
}

// BlendTransition is the style a transition draws at a moment in time: the
// target [Style] is blended in from the entry [Style] for every transition
// it declared. elapsed is milliseconds since the transition began. While a
// transition waits out its delay the entry's value shows; past its end the
// target shows, so a finished transition is the target itself.
func BlendTransition(from, to Style, elapsed float64, ctx Units) Style {
	out := CopyStyle(to)
	for _, tx := range to.Transitions {
		start := tx.Delay.MS()
		dur := tx.Duration.MS()
		if dur <= 0 {
			continue
		}
		if tx.Prop == "all" {
			for _, prop := range animatableProps {
				blendPropTimed(&out, from, prop, start, dur, elapsed, tx.Timing, ctx)
			}
			continue
		}
		if tx.Prop == "none" {
			continue
		}
		blendPropTimed(&out, from, tx.Prop, start, dur, elapsed, tx.Timing, ctx)
	}
	return out
}

// BlendKeyframes folds the two frames straddling a point in time into the
// out style: for every animated property the frames wrote, the value at the
// progress between them. A property one frame omitted takes the base style's
// value as its other endpoint, so a frame that sets only half a transform
// list still travels from the base.
func BlendKeyframes(out *Style, base Style, fa, fb Style, props map[string]bool, local float64, ctx Units) {
	fromV := CopyStyle(base)
	toV := CopyStyle(base)
	for prop := range props {
		if fa.Set[prop] {
			setProp(&fromV, fa, prop)
		}
		if fb.Set[prop] {
			setProp(&toV, fb, prop)
		}
		blendProp(&toV, fromV, prop, local, ctx)
		setProp(out, toV, prop)
	}
}

// overlayProps folds every property in props from src into dst, marking each
// in dst's Set map.
func overlayProps(dst *Style, src Style, props map[string]bool) {
	for prop := range props {
		setProp(dst, src, prop)
	}
}

// blendPropTimed runs one property through a whole transition: the entry's
// value until the delay passes, the target's past the duration, and the
// easing between. out holds the target already, so only the wait and the
// middle need writing.
func blendPropTimed(out *Style, from Style, prop string, start, dur, elapsed float64, tim Timing, ctx Units) {
	if elapsed <= start {
		setProp(out, from, prop)
		return
	}
	local := (elapsed - start) / dur
	if local >= 1 {
		return
	}
	blendProp(out, from, prop, tim.Ease(local), ctx)
}

// setProp copies one canonical property's value from src to dst and marks it
// in dst's Set map.
func setProp(dst *Style, src Style, prop string) {
	switch prop {
	case "background-color":
		dst.Background = src.Background
	case "color":
		dst.Color = src.Color
	case "opacity":
		dst.Opacity = src.Opacity
	case "width":
		dst.Width = src.Width
	case "height":
		dst.Height = src.Height
	case "min-width":
		dst.MinWidth = src.MinWidth
	case "max-width":
		dst.MaxWidth = src.MaxWidth
	case "min-height":
		dst.MinHeight = src.MinHeight
	case "max-height":
		dst.MaxHeight = src.MaxHeight
	case "margin-top":
		dst.Margin[0] = src.Margin[0]
	case "margin-right":
		dst.Margin[1] = src.Margin[1]
	case "margin-bottom":
		dst.Margin[2] = src.Margin[2]
	case "margin-left":
		dst.Margin[3] = src.Margin[3]
	case "padding-top":
		dst.Padding[0] = src.Padding[0]
	case "padding-right":
		dst.Padding[1] = src.Padding[1]
	case "padding-bottom":
		dst.Padding[2] = src.Padding[2]
	case "padding-left":
		dst.Padding[3] = src.Padding[3]
	case "top":
		dst.Top = src.Top
	case "right":
		dst.Right = src.Right
	case "bottom":
		dst.Bottom = src.Bottom
	case "left":
		dst.Left = src.Left
	case "font-size":
		dst.FontSize = src.FontSize
	case "letter-spacing":
		dst.LetterSpacing = src.LetterSpacing
	case "word-spacing":
		dst.WordSpacing = src.WordSpacing
	case "line-height":
		dst.LineHeight = src.LineHeight
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
	if dst.Set != nil {
		dst.Set[prop] = true
	}
}

// blendProp interpolates one canonical property from a to b, writing the
// result into b (and its Set map): dst starts as the target value, which is
// exactly how a blend composes. Discrete properties swap over at the half
// way point.
func blendProp(dst *Style, from Style, prop string, t float64, ctx Units) {
	switch prop {
	case "background-color":
		dst.Background = canvas.Mix(from.Background, dst.Background, float32(t))
	case "color":
		dst.Color = canvas.Mix(from.Color, dst.Color, float32(t))
	case "opacity":
		dst.Opacity = from.Opacity + (dst.Opacity-from.Opacity)*t
	case "width":
		dst.Width = blendLength(from.Width, dst.Width, t, ctx)
	case "height":
		dst.Height = blendLength(from.Height, dst.Height, t, ctx)
	case "min-width":
		dst.MinWidth = blendLength(from.MinWidth, dst.MinWidth, t, ctx)
	case "max-width":
		dst.MaxWidth = blendLength(from.MaxWidth, dst.MaxWidth, t, ctx)
	case "min-height":
		dst.MinHeight = blendLength(from.MinHeight, dst.MinHeight, t, ctx)
	case "max-height":
		dst.MaxHeight = blendLength(from.MaxHeight, dst.MaxHeight, t, ctx)
	case "margin-top":
		dst.Margin[0] = blendLength(from.Margin[0], dst.Margin[0], t, ctx)
	case "margin-right":
		dst.Margin[1] = blendLength(from.Margin[1], dst.Margin[1], t, ctx)
	case "margin-bottom":
		dst.Margin[2] = blendLength(from.Margin[2], dst.Margin[2], t, ctx)
	case "margin-left":
		dst.Margin[3] = blendLength(from.Margin[3], dst.Margin[3], t, ctx)
	case "padding-top":
		dst.Padding[0] = blendLength(from.Padding[0], dst.Padding[0], t, ctx)
	case "padding-right":
		dst.Padding[1] = blendLength(from.Padding[1], dst.Padding[1], t, ctx)
	case "padding-bottom":
		dst.Padding[2] = blendLength(from.Padding[2], dst.Padding[2], t, ctx)
	case "padding-left":
		dst.Padding[3] = blendLength(from.Padding[3], dst.Padding[3], t, ctx)
	case "top":
		dst.Top = blendLength(from.Top, dst.Top, t, ctx)
	case "right":
		dst.Right = blendLength(from.Right, dst.Right, t, ctx)
	case "bottom":
		dst.Bottom = blendLength(from.Bottom, dst.Bottom, t, ctx)
	case "left":
		dst.Left = blendLength(from.Left, dst.Left, t, ctx)
	case "font-size":
		dst.FontSize = int(float64(from.FontSize) + float64(dst.FontSize-from.FontSize)*t + 0.5)
	case "letter-spacing":
		dst.LetterSpacing = blendLength(from.LetterSpacing, dst.LetterSpacing, t, ctx)
	case "word-spacing":
		dst.WordSpacing = blendLength(from.WordSpacing, dst.WordSpacing, t, ctx)
	case "line-height":
		dst.LineHeight = from.LineHeight + (dst.LineHeight-from.LineHeight)*t
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
	if dst.Set != nil {
		dst.Set[prop] = true
	}
}

// blendLength interpolates two lengths. Same-unit lengths lerp their values,
// so a width that goes 50% to 75% stays a percentage all the way. Keywords
// and units that cannot mix resolve to their reference pixels and lerp those.
func blendLength(a, b Length, t float64, ctx Units) Length {
	if a.u == b.u && a.u != unitAuto && a.u != unitNone {
		return Length{u: a.u, value: a.value + (b.value-a.value)*t}
	}
	if a.u == unitAuto || a.u == unitNone || b.u == unitAuto || b.u == unitNone {
		if t < 0.5 {
			return a
		}
		return b
	}
	ar, br := a.Resolve(ctx), b.Resolve(ctx)
	return Fixed(float64(ar) + float64(br-ar)*t)
}

// blendTransformList blends two transform lists, function by function. Lists
// of the same length with matching kinds interpolate; anything else swaps
// wholesale at the half way point, which is how CSS treats a transform that
// cannot match up.
func blendTransformList(a, b []TransformFunc, t float64) []TransformFunc {
	if len(a) != len(b) {
		if t < 0.5 {
			return a
		}
		return b
	}
	if len(a) == 0 {
		return nil
	}
	out := make([]TransformFunc, 0, len(a))
	for i := range a {
		fa, fb := a[i], b[i]
		if fa.Kind != fb.Kind {
			if t < 0.5 {
				return a
			}
			return b
		}
		combine := blendTransformFunc(fa, fb, t)
		out = append(out, combine)
	}
	return out
}

func blendTransformFunc(fa, fb TransformFunc, t float64) TransformFunc {
	out := fa
	switch fa.Kind {
	case TransformTranslate:
		out.Dx = blendLength(fa.Dx, fb.Dx, t, Units{})
		out.Dy = blendLength(fa.Dy, fb.Dy, t, Units{})
	case TransformScale:
		out.Sx = fa.Sx + (fb.Sx-fa.Sx)*t
		out.Sy = fa.Sy + (fb.Sy-fa.Sy)*t
	case TransformRotate:
		out.Ax = Rad(fa.Ax.Rad() + (fb.Ax.Rad()-fa.Ax.Rad())*t)
	case TransformSkew:
		out.Ax = Rad(fa.Ax.Rad() + (fb.Ax.Rad()-fa.Ax.Rad())*t)
		out.Ay = Rad(fa.Ay.Rad() + (fb.Ay.Rad()-fa.Ay.Rad())*t)
	case TransformMatrix:
		out.M = [6]float64{}
		for j := 0; j < 6; j++ {
			out.M[j] = fa.M[j] + (fb.M[j]-fa.M[j])*t
		}
	}
	return out
}

// blendShadowList blends two shadow lists, shadow by shadow. Lists of equal
// length blend in step; a list that grew or shrank swaps at the half way
// point, since there is no honest better guess.
func blendShadowList(a, b []Shadow, t float64) []Shadow {
	if len(a) != len(b) {
		if t < 0.5 {
			return a
		}
		return b
	}
	if len(a) == 0 {
		return nil
	}
	out := make([]Shadow, len(a))
	for i := range a {
		out[i] = blendShadow(a[i], b[i], t)
	}
	return out
}

func blendShadow(a, b Shadow, t float64) Shadow {
	if a.Inset != b.Inset {
		if t < 0.5 {
			return a
		}
		return b
	}
	return Shadow{
		X:      roundLerp(a.X, b.X, t),
		Y:      roundLerp(a.Y, b.Y, t),
		Blur:   roundLerp(a.Blur, b.Blur, t),
		Spread: roundLerp(a.Spread, b.Spread, t),
		Color:  canvas.Mix(a.Color, b.Color, float32(t)),
		Inset:  a.Inset,
	}
}

// blendFilterList blends two filter lists, filter by filter, with the same
// length rule as the shadows. A drop-shadow pair blends; a pair of the same
// function kind lerps its amount; anything else swaps at the half way point.
func blendFilterList(a, b []Filter, t float64) []Filter {
	if len(a) != len(b) {
		if t < 0.5 {
			return a
		}
		return b
	}
	if len(a) == 0 {
		return nil
	}
	out := make([]Filter, len(a))
	for i := range a {
		out[i] = blendFilter(a[i], b[i], t)
	}
	return out
}

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

func roundLerp(a, b int, t float64) int {
	return int(float64(a) + float64(b-a)*t + 0.5)
}
