package template

import "github.com/gabrielluizsf/antui"

// dateGeom is the picker's box and the calendar that opens over the window.
type dateGeom struct {
	u              int // the scale the calendar draws at, not the window's
	x, y, w, h     int // the box
	px, py, dw, dh int // the popup, centred in the window it is on
}

// dateGeomFor measures the picker: the box where the field stands, and the
// calendar grown by the biggest whole step that still fits the window and
// centred in it, so an open calendar is the screen rather than a slip of paper
// hanging off the field. Whole steps are what the calendar is measured in —
// every part of it is one number times this scale — so growing it by a step
// keeps the days, the text and the arrows in step with each other.
func dateGeomFor(win *antui.Window, x, y, w, h int) dateGeom {
	u := Scale(win)
	bw, bh := datePickerSize(u)
	steps := min((win.Width()*9)/10/bw, (win.Height()*8)/10/bh)
	pu := u * max(steps, 1)
	dw, dh := datePickerSize(pu)
	return dateGeom{
		u: pu, x: x, y: y, w: w, h: h,
		px: (win.Width() - dw) / 2,
		py: (win.Height() - dh) / 2,
		dw: dw, dh: dh,
	}
}
