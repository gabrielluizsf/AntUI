package template

import "github.com/gabrielluizsf/antui/canvas"

// paintCyberDateDays draws the day cells, marking the selected, hovered and
// today.
func paintCyberDateDays(cv *canvas.Canvas, v vivid, u, x, y, firstWD, days, selected, today, hover int) {
	for day := 1; day <= days; day++ {
		pos := firstWD + day - 1
		col, row := pos%7, pos/7
		cx, cy, cc, _ := dayCellRect(x, y, u, col, row)
		switch day {
		case selected:
			glow(cv, cx, cy, cc, cc, 2*u, v.neon, int(float64(20)*v.glow), u)
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, v.neon)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), cyberPalette.bg, u)
		case hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, cyberPalette.panelHov)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), v.text, u)
		default:
			colour := v.text
			if day == today {
				colour = v.neon
			}
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), colour, u)
		}
		if day == today && day != selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2-2*u, v.neon)
		}
	}
}
