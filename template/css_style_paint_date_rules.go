package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintDateRules draws the fine grid the body is read as a table by: its rules
// make the weekday columns and week rows easy to scan while staying underneath
// the selected and hovered days.
func (cs *cssStyle) paintDateRules(cv *canvas.Canvas, st css.Style, p datePaint) {
	u := p.u
	gx, gy, gw, _ := dateGridRect(p.x, p.y, u)
	c := dateCellSize(u)
	rule := cs.fade(st, canvas.Shade(cs.fill(st), -0.035))
	for col := 0; col <= 7; col++ {
		cv.FillRect(canvas.Area{X: gx + col*c, Y: p.weekHeadY(), Width: u, Height: weekdayHeadHeight(u) + p.rows()*c}, rule)
	}
	for row := 0; row <= p.rows(); row++ {
		cv.FillRect(canvas.Area{X: gx, Y: gy + row*c, Width: gw, Height: u}, rule)
	}
}
