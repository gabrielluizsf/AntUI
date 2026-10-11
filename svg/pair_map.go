package svg

import "github.com/gabrielluizsf/antui/canvas"

// pairMap runs a two-input primitive's colour over every pixel the two
// pictures have, the one it is given first and the one it takes second, and
// writes the answer into the first. A pixel of one picture the other does not
// reach keeps what it had.
func pairMap(dst, src *canvas.Canvas, f func(a, b canvas.Color) canvas.Color) {
	if dst == nil || src == nil || dst.Pixels == nil || src.Pixels == nil {
		return
	}
	for y := 0; y < dst.Height && y < src.Height; y++ {
		drow := dst.Pixels[y*dst.Stride : y*dst.Stride+dst.Width]
		srow := src.Pixels[y*src.Stride : y*src.Stride+src.Width]
		for x := 0; x < dst.Width && x < src.Width; x++ {
			drow[x] = f(drow[x], srow[x])
		}
	}
}
