package template

import "github.com/gabrielluizsf/antui/canvas"

// paintCyberDateDays draws the day cells the way the site's calendar does:
// a whisper-white cell, hover lifts the edge to acc-dim, selected washes
// the cell in accent behind an acc border and a small glow, today gets a
// half-strength accent frame.
func paintCyberDateDays(cv *canvas.Canvas, v vivid, p datePaint) {
	for day := 1; day <= p.days; day++ {
		pos := p.firstWD + day - 1
		cx, cy, cc, _ := dayCellRect(p.x, p.y, p.u, pos%7, pos/7)
		num := itoa(day)
		tx, ty := cx+(cc-textWidth(p.u, num))/2, cy+(cc-textHeight(p.u))/2
		cell := canvas.Round(canvas.Area{X: cx, Y: cy, Width: cc, Height: cc}, 2*p.u)
		switch day {
		case p.selected:
			glow(cv, halo{x: cx, y: cy, w: cc, h: cc, radius: 2 * p.u, colour: v.neon, strength: int(float64(12) * v.glow), u: p.u})
			cv.FillRoundRect(cell, canvas.Fade(v.neon, 51))
			cv.RoundRect(cell, v.neon)
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: v.text, Scale: p.u})
		case p.hover:
			cv.FillRoundRect(cell, canvas.Fade(canvas.White, 5))
			cv.RoundRect(cell, canvas.Fade(v.neon, 110))
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: v.text, Scale: p.u})
		default:
			cv.FillRoundRect(cell, canvas.Fade(canvas.White, 5))
			colour := v.text
			if day == p.today {
				colour = v.neon
			}
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: colour, Scale: p.u})
		}
		if day == p.today && day != p.selected {
			cv.RoundRect(cell, canvas.Fade(v.neon, 128))
		}
	}
}
