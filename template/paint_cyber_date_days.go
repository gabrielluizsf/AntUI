package template

import "github.com/gabrielluizsf/antui/canvas"

// paintCyberDateDays draws the day cells, marking the selected, hovered and
// today.
func paintCyberDateDays(cv *canvas.Canvas, v vivid, p datePaint) {
	for day := 1; day <= p.days; day++ {
		pos := p.firstWD + day - 1
		cx, cy, cc, _ := dayCellRect(p.x, p.y, p.u, pos%7, pos/7)
		num := itoa(day)
		tx, ty := cx+(cc-textWidth(p.u, num))/2, cy+(cc-textHeight(p.u))/2
		switch day {
		case p.selected:
			glow(cv, halo{x: cx, y: cy, w: cc, h: cc, radius: 2 * p.u, colour: v.neon, strength: int(float64(20) * v.glow), u: p.u})
			cv.FillRoundRect(cx, cy, cc, cc, 2*p.u, v.neon)
			cv.TextScaled(tx, ty, num, cyberPalette.bg, p.u)
		case p.hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*p.u, cyberPalette.panelHov)
			cv.TextScaled(tx, ty, num, v.text, p.u)
		default:
			colour := v.text
			if day == p.today {
				colour = v.neon
			}
			cv.TextScaled(tx, ty, num, colour, p.u)
		}
		if day == p.today && day != p.selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2-2*p.u, v.neon)
		}
	}
}
