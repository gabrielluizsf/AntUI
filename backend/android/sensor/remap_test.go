package sensor

import "testing"

// A phone lying flat, tipped so that the device's own X points right and
// down the screen. What the screen calls "right" depends on how far the
// screen has been turned, and the sensor does not turn with it.
func TestRemapFollowsTheScreen(t *testing.T) {
	// A push along the device's X axis, and nothing along Y.
	const push = 9.8
	for _, c := range []struct {
		rotation  int
		wantX     float32
		wantY     float32
		looksLike string
	}{
		{0, push, 0, "the same way: the screen is as the device was built"},
		{1, 0, push, "up the screen: the device is on its side"},
		{2, -push, 0, "the other way: the device is upside down"},
		{3, 0, -push, "down the screen"},
	} {
		x, y := remap(push, 0, c.rotation)
		if x != c.wantX || y != c.wantY {
			t.Errorf("rotation %d: (%g, %g), want (%g, %g) — %s",
				c.rotation, x, y, c.wantX, c.wantY, c.looksLike)
		}
	}
}

// Four quarter turns is where it started, whatever the vector.
func TestRemapComesBackRound(t *testing.T) {
	x, y := float32(3), float32(-7)
	a, b := x, y
	for range 4 {
		a, b = remap(a, b, 1)
	}
	if a != x || b != y {
		t.Errorf("four quarter turns gave (%g, %g), want (%g, %g)", a, b, x, y)
	}
}

// A rotation outside 0..3 is the same as its remainder, so a caller that
// passes degrees by mistake gets something wrong rather than something
// undefined — and a caller that passes 4 gets 0.
func TestRemapWrapsRatherThanIndexing(t *testing.T) {
	for _, r := range []int{4, 8, -4} {
		if x, y := remap(1, 2, r); x != 1 || y != 2 {
			t.Errorf("rotation %d gave (%g, %g), want it treated as 0", r, x, y)
		}
	}
}
