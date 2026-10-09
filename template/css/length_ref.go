package css

// ref converts the length into reference pixels under a context, the unit
// arithmetic every resolution is built on. Percentages are of the containing
// width, viewport units of the window, font units of the font sizes; fixed and
// physical lengths are their own number of pixels.
func (l Length) ref(ctx Units) float64 {
	switch l.u {
	case unitPx:
		return l.value
	case unitPct:
		return l.value * float64(ctx.Width) / 100
	case unitEm:
		return l.value * float64(ctx.font())
	case unitRem:
		return l.value * float64(ctx.root())
	case unitVw:
		return l.value * float64(ctx.Width) / 100
	case unitVh:
		return l.value * float64(ctx.Height) / 100
	case unitVmin:
		return l.value * float64(min(ctx.Width, ctx.Height)) / 100
	case unitVmax:
		return l.value * float64(max(ctx.Width, ctx.Height)) / 100
	// The advance of "0" and the x-height both settle around half an em for
	// the built-in face; the engine reads both as 0.5em.
	case unitCh:
		return l.value * 0.5 * float64(ctx.font())
	case unitEx:
		return l.value * 0.5 * float64(ctx.font())
	case unitCm:
		return l.value * pxPerCm
	case unitMm:
		return l.value * pxPerMm
	case unitIn:
		return l.value * pxPerIn
	case unitPt:
		return l.value * pxPerPt
	case unitPc:
		return l.value * pxPerPc
	case unitQ:
		return l.value * pxPerQ
	}
	return 0
}
