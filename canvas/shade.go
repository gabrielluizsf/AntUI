package canvas

// Shade lightens the colour for a positive amount and darkens it for a
// negative one, from -1 to 1. The alpha is left alone.
func Shade(c Color, amount float32) Color {
	var target Color
	if amount >= 0 {
		target = (c & 0xFF000000) | 0x00FFFFFF
	} else {
		target = c & 0xFF000000
		amount = -amount
	}
	return Mix(c, target, amount)
}
