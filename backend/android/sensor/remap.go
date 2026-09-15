package sensor

// The display remap is arithmetic and nothing else, so it has no build tag
// and is tested on the machine that builds rather than on a phone. It is
// also the single most common way to get sensors wrong, which is reason
// enough to be able to test it at all.

// remap turns a vector from the device's fixed frame into the screen's, for
// a screen rotated by rotation quarter-turns clockwise.
func remap(x, y float32, rotation int) (float32, float32) {
	switch rotation & 3 {
	case 1: // a quarter turn
		return -y, x
	case 2: // upside down
		return -x, -y
	case 3: // three quarters
		return y, -x
	}
	return x, y
}
