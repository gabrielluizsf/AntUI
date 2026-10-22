package template

import "github.com/gabrielluizsf/antui/canvas"

// paintCyberMonthGrid draws the twelve months of the year a tap on the
// header's title opened, the chosen one lit and glowing the way a chosen day
// is.
func paintCyberMonthGrid(cv *canvas.Canvas, v vivid, p datePaint) {
	inset := 3 * p.u
	for m := 1; m <= 12; m++ {
		col, row := (m-1)%4, (m-1)/4
		cx, cy, cw, ch := monthCellRect(p.x, p.y, p.u, col, row)
		bx, by, bw, bh := cx+inset, cy+inset, cw-2*inset, ch-2*inset
		name := monthAbbr(m)
		tx := bx + (bw-textWidth(p.u, name))/2
		ty := by + (bh-textHeight(p.u))/2
		if m == p.month {
			cell := canvas.Round(canvas.Area{X: bx, Y: by, Width: bw, Height: bh}, 2*p.u)
			glow(cv, halo{x: bx, y: by, w: bw, h: bh, radius: 2 * p.u, colour: v.neon, strength: int(float64(12) * v.glow), u: p.u})
			cv.FillRoundRect(cell, canvas.Fade(v.neon, 51))
			cv.RoundRect(cell, v.neon)
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, name, canvas.TextStyle{Color: v.text, Scale: p.u})
		} else {
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, name, canvas.TextStyle{Color: v.text, Scale: p.u})
		}
	}
}
