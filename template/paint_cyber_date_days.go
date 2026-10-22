package template

import "github.com/gabrielluizsf/antui/canvas"

// paintCyberDateDays draws the day cells the way the site's calendar does: a
// chamfered mark a couple of units short of its cell on every side — the gap
// is what separates one day from the next — hover lifts the edge to acc-dim,
// selected washes it in accent behind an acc border and a small glow, today
// gets a half-strength accent frame. An ordinary day carries no mark at all;
// the ones that do are what the eye should find first.
func paintCyberDateDays(cv *canvas.Canvas, v vivid, p datePaint) {
	inset := 2 * p.u
	for day := 1; day <= p.days; day++ {
		pos := p.firstWD + day - 1
		cx, cy, cc, _ := dayCellRect(p.x, p.y, p.u, pos%7, pos/7)
		bx, by, bs := cx+inset, cy+inset, cc-2*inset
		num := itoa(day)
		tx, ty := cx+(cc-textWidth(p.u, num))/2, cy+(cc-textHeight(p.u))/2
		cell := canvas.Round(canvas.Area{X: bx, Y: by, Width: bs, Height: bs}, 2*p.u)
		switch day {
		case p.selected:
			glow(cv, halo{x: bx, y: by, w: bs, h: bs, radius: 2 * p.u, colour: v.neon, strength: int(float64(12) * v.glow), u: p.u})
			cv.FillRoundRect(cell, canvas.Fade(v.neon, 51))
			cv.RoundRect(cell, v.neon)
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: v.text, Scale: p.u})
		case p.hover:
			cv.FillRoundRect(cell, canvas.Fade(canvas.White, 10))
			cv.RoundRect(cell, canvas.Fade(v.neon, 110))
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: v.text, Scale: p.u})
		default:
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
