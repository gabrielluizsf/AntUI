package svg

import "github.com/gabrielluizsf/antui/canvas"

// contextFill is what a `context-fill` inside a marker resolves to: the single
// colour the element that named the marker fills with. A fill that is a
// gradient or a pattern is not one colour and is not carried across; a fill that
// is itself a context paint carries the one already in hand.
//
// shortcut: gradient and pattern fills do not reach context-fill, upgrade when a
// marker has to paint with the naming shape's paint server.
func (st Style) contextFill(paints palette) canvas.Color {
	switch {
	case !st.HasFill || st.fillGradient != nil || st.fillPattern != nil:
		return 0
	case st.FillCurrent:
		return paints.current
	case st.FillContext:
		return paints.fill
	default:
		return st.Fill
	}
}

// contextStroke is what a `context-stroke` inside a marker resolves to: the
// single colour the element that named the marker strokes with, read and
// limited the same way [Style.contextFill] reads a fill.
func (st Style) contextStroke(paints palette) canvas.Color {
	switch {
	case !st.HasStroke || st.strokeGradient != nil || st.strokePattern != nil:
		return 0
	case st.StrokeCurrent:
		return paints.current
	case st.StrokeContext:
		return paints.stroke
	default:
		return st.Stroke
	}
}
