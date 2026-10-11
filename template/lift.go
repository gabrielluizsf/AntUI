package template

import "github.com/gabrielluizsf/antui/canvas"

// lift brings a colour partway to white. The weight is (f-1)/f of the way:
// identity at f = 1, halfway at 50% opacity, capped two thirds of the way
// there. RGB climbs, alpha stays.
func lift(c canvas.Color, f float64) canvas.Color {
	if f <= 1 {
		return c
	}
	alpha := byte(c >> 24)
	w := int((f - 1) * 255 / f)
	up := func(v byte) byte {
		return byte(int(v) + (255-int(v))*w/255)
	}
	return canvas.Color(alpha)<<24 |
		canvas.Color(up(byte(c>>16)))<<16 |
		canvas.Color(up(byte(c>>8)))<<8 |
		canvas.Color(up(byte(c)))
}
