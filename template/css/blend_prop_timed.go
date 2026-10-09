package css

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
