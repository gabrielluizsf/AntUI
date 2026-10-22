package template

import "github.com/gabrielluizsf/antui"

// dateGeom is the picker's box and the calendar that opens over the window.
type dateGeom struct {
	u              int // the scale the calendar draws at, not the window's
	x, y, w, h     int // the box
	px, py, dw, dh int // the popup, centred in the window it is on
}

// dateMaxExtraSteps is the most a calendar is ever grown past the rest of
// the window's own scale. Filling the window by whatever whole step
// happens to fit once drew the day numbers at three or four times
// everything else on screen, which read as a jump in scale rather than a
// popup standing out — a dialog earns its emphasis from sitting on a
// darkened scrim, not from its own text towering over the field it opened
// from.
const dateMaxExtraSteps = 1

// dateGeomFor measures the picker: the box where the field stands, and the
// calendar grown by the biggest whole step that still fits the window and
// centred in it, capped at [dateMaxExtraSteps] past the window's own scale.
// Whole steps are what the calendar is measured in — every part of it is
// one number times this scale — so growing it by a step keeps the days, the
// text and the arrows in step with each other.
func dateGeomFor(win *antui.Window, x, y, w, h int) dateGeom {
	u := Scale(win)
	bw, bh := datePickerSize(u)
	steps := min((win.Width()*9)/10/bw, (win.Height()*8)/10/bh, 1+dateMaxExtraSteps)
	pu := u * max(steps, 1)
	dw, dh := datePickerSize(pu)
	return dateGeom{
		u: pu, x: x, y: y, w: w, h: h,
		px: (win.Width() - dw) / 2,
		py: (win.Height() - dh) / 2,
		dw: dw, dh: dh,
	}
}
