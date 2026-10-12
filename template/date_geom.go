package template

// dateGeom is the picker's box and the calendar that pops up below it.
type dateGeom struct {
	u              int
	x, y, w, h     int // the box
	px, py, dw, dh int // the popup
}

// dateGeomFor measures the picker at scale u.
func dateGeomFor(x, y, w, h, u int) dateGeom {
	dw, dh := datePickerSize(u)
	return dateGeom{u: u, x: x, y: y, w: w, h: h, px: x, py: y + h, dw: dw, dh: dh}
}
