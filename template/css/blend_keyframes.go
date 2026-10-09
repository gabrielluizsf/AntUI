package css

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
