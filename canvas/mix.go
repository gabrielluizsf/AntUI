package canvas

// Mix interpolates between two colours, t running from 0 (all a) to 1 (all
// b). Every channel is interpolated, alpha included.
func Mix(a, b Color, t float32) Color {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	k := uint32(t*255 + 0.5)
	inv := 255 - k
	ch := func(shift uint) uint32 {
		return ((uint32(a>>shift)&0xFF)*inv + (uint32(b>>shift)&0xFF)*k) / 255
	}
	return Color(ch(24)<<24 | ch(16)<<16 | ch(8)<<8 | ch(0))
}
