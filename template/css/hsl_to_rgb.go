package css

// hslToRGB converts hue degrees and saturation/lightness percentages to sRGB.
func hslToRGB(h, s, l float64) [3]uint8 {
	h = h / 360
	s = s / 100
	l = l / 100
	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		q := l + s - l*s
		if l < 0.5 {
			q = l * (1 + s)
		}
		p := 2*l - q
		hu := func(t float64) float64 {
			if t < 0 {
				t++
			}
			if t > 1 {
				t--
			}
			switch {
			case t < 1.0/6.0:
				return p + (q-p)*6*t
			case t < 1.0/2.0:
				return q
			case t < 2.0/3.0:
				return p + (q-p)*(2.0/3.0-t)*6
			}
			return p
		}
		r = hu(h + 1.0/3.0)
		g = hu(h)
		b = hu(h - 1.0/3.0)
	}
	return [3]uint8{uint8(r*255 + 0.5), uint8(g*255 + 0.5), uint8(b*255 + 0.5)}
}
