package svg

import "github.com/gabrielluizsf/antui/canvas"

// colorMatrix runs feColorMatrix over a picture. The matrix works on the
// colour as it stands — not on it pressed against the background — so the
// four channels of every pixel go in straight and the four that come out are
// clamped and written back.
func colorMatrix(cv *canvas.Canvas, m [20]float64) {
	if cv == nil || cv.Pixels == nil {
		return
	}
	for y := 0; y < cv.Height; y++ {
		row := cv.Pixels[y*cv.Stride : y*cv.Stride+cv.Width]
		for x, c := range row {
			r, g, b, a := frac(c.R()), frac(c.G()), frac(c.B()), frac(c.A())
			rp := clamp01(m[0]*r + m[1]*g + m[2]*b + m[3]*a + m[4])
			gp := clamp01(m[5]*r + m[6]*g + m[7]*b + m[8]*a + m[9])
			bp := clamp01(m[10]*r + m[11]*g + m[12]*b + m[13]*a + m[14])
			ap := clamp01(m[15]*r + m[16]*g + m[17]*b + m[18]*a + m[19])
			row[x] = rgbaColor(rp, gp, bp, ap)
		}
	}
}
