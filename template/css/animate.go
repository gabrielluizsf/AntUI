package css

// animate applies a @keyframes block at the progress t (0..1) over the base
// style and returns the animated style. The result owns its Set map, so the
// cached base style is never touched; properties the frames wrote overlay
// the base, and between two frames every animatable property the frames
// mentioned interpolates, with the base style's value standing in for a
// frame that omitted it. Properties the engine does not interpolate jump at
// the 0.5 crossing.
type AnimateParams struct {
	Base      Style
	Keyframes *Keyframes
	Progress  float64
	Units     Units
}

func Animate(p AnimateParams) Style {
	base, kf, t, ctx := p.Base, p.Keyframes, p.Progress, p.Units
	if kf == nil || len(kf.Frames) == 0 {
		return base
	}
	out := CopyStyle(base)
	lo, hi, local, ok := frameSpan(kf.Frames, t)
	if !ok {
		return out
	}
	if lo.Offset == hi.Offset {
		fa := frameStyle(lo, base, ctx)
		overlayProps(&out, fa, fa.Set)
		return out
	}
	fa := frameStyle(lo, base, ctx)
	fb := frameStyle(hi, base, ctx)
	props := map[string]bool{}
	for _, prop := range animatableProps {
		if fa.Set[prop] || fb.Set[prop] {
			props[prop] = true
		}
	}
	blendKeyframes(&out, base, fa, fb, props, local, ctx)
	return out
}

func animate(base Style, kf *Keyframes, t float64, ctx Units) Style {
	return Animate(AnimateParams{Base: base, Keyframes: kf, Progress: t, Units: ctx})
}
