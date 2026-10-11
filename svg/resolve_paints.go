package svg

// resolvePaints follows the `url(#id)` a fill or a stroke named, which only the
// whole drawing can say anything about. What the reference turns into is the
// gradient or the pattern that was in the drawing; the colour written after the
// `url(...)` where there was not; and nothing, with one warning, where the
// shape asked for something that is not there and gave no second choice.
//
// The reference is cleared as it is spent, so that a shape inside a group with a
// missing gradient does not say the same thing once for every shape inside it.
func (img *Image) resolvePaints(st *Style, warn func(string, ...any)) {
	if st.fillRef != "" {
		g, p := img.grads[st.fillRef], img.patterns[st.fillRef]
		st.fillRef, st.fillGradient, st.fillPattern = "", g, p
		switch {
		case g != nil || p != nil:
			st.HasFill = true
		case st.hasFillFallback:
			st.Fill, st.HasFill = st.fillFallback, true
			st.FillCurrent = st.fillFallbackCurrent
		default:
			st.HasFill = false
			warn("the fill names a gradient or a pattern that is not in the drawing and has no colour after it, so the shape is not filled")
		}
	}
	if st.strokeRef != "" {
		g, p := img.grads[st.strokeRef], img.patterns[st.strokeRef]
		st.strokeRef, st.strokeGradient, st.strokePattern = "", g, p
		switch {
		case g != nil || p != nil:
			st.HasStroke = true
		case st.hasStrokeFallback:
			st.Stroke, st.HasStroke = st.strokeFallback, true
			st.StrokeCurrent = st.strokeFallbackCurrent
		default:
			st.HasStroke = false
			warn("the stroke names a gradient or a pattern that is not in the drawing and has no colour after it, so the shape is not outlined")
		}
	}
}
