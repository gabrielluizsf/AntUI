package template

// The steps a widget's own surface colour is shaded by to paint the parts that
// are not the widget: the empty track under a value, the cell under the
// pointer, the thumb sitting on its own track. A negative number darkens, and
// the number is how far the colour moves toward black — the same convention
// [canvas.Shade] takes.
const (
	// trackShade is the empty part of a control that holds a value: the
	// groove of a slider, the bar behind a progress, the full ring the
	// value's arc sweeps over.
	trackShade = -0.25

	// hoverFillShade is a cell under the pointer — a dropdown's option, a
	// day in the calendar. The lightest of the three, because it sits
	// beside other cells and a stack of them has to stay a list.
	hoverFillShade = -0.15

	// thumbShade is a switch's thumb, darker than a hovered cell because
	// there is only one of it and it is the thing being looked at.
	thumbShade = -0.20

	// hoverBorderShade is a control's border when it is hovered. The
	// darkest, because a one-pixel line carries almost no weight of its own
	// and has to move further to be noticed.
	hoverBorderShade = -0.30
)

// trackAlpha is how opaque a scrollbar's trough is over the surface behind it:
// a half, which is enough to see where the thumb can travel without the trough
// reading as a hole cut in the page.
const trackAlpha = 128
