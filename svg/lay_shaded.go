package svg

import "github.com/gabrielluizsf/antui/canvas"

// layShaded paints the writing where the colour comes from where each pixel
// lands rather than being one colour over all of it, which is what a gradient
// and a pattern both are. The letters go into a picture of their own first, so
// that there is a coverage per pixel to ask about, and the colour is then laid
// over what is already there one pixel at a time — the same question a shape
// filled with either asks of every pixel it covers, through FillPathFunc.
func layShaded(cv *canvas.Canvas, put putRun, shade func(x, y int) canvas.Color) {
	layer, err := canvas.NewLayer(cv.Width, cv.Height)
	if err != nil {
		return
	}
	// Any opaque colour will do: only what it covered matters, and the colour
	// comes from the shade at the pixel rather than from this.
	put(layer, canvas.White, 0, 0)
	for y := range layer.Height {
		for x := range layer.Width {
			a := layer.At(x, y).A()
			if a == 0 {
				continue
			}
			c := shade(x, y)
			if c.A() == 0 {
				continue
			}
			if a < 255 {
				c = canvas.RGBA(c.R(), c.G(), c.B(), uint8(int(c.A())*int(a)/255))
			}
			cv.Pixel(x, y, c)
		}
	}
}
