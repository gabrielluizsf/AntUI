package canvas

import "testing"

func TestRGBChannels(t *testing.T) {
	c := RGB(0x11, 0x22, 0x33)
	if c != 0xFF112233 {
		t.Errorf("RGB = %08X, want FF112233", uint32(c))
	}
	if got := [4]uint8{c.R(), c.G(), c.B(), c.A()}; got != [4]uint8{0x11, 0x22, 0x33, 0xFF} {
		t.Errorf("channels = %v, want [11 22 33 FF]", got)
	}
	if got := RGBA(1, 2, 3, 4); got != 0x04010203 {
		t.Errorf("RGBA = %08X, want 04010203", uint32(got))
	}
}

func TestBlend(t *testing.T) {
	tests := []struct {
		name     string
		dst, src Color
		want     Color
	}{
		{"opaque source replaces the background", White, Black, Black},
		{"transparent source keeps the background", White, Transparent, White},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Blend(tt.dst, tt.src); got != tt.want {
				t.Errorf("Blend = %08X, want %08X", uint32(got), uint32(tt.want))
			}
		})
	}

	half := Blend(Black, Fade(White, 128))
	if half.R() < 120 || half.R() > 135 {
		t.Errorf("alpha 128 over black gave red %d, want it near halfway", half.R())
	}
}

func TestMix(t *testing.T) {
	if got := Mix(Black, White, 0); got != Black {
		t.Errorf("Mix(t=0) = %08X, want the first colour", uint32(got))
	}
	if got := Mix(Black, White, 1); got != White {
		t.Errorf("Mix(t=1) = %08X, want the second colour", uint32(got))
	}
	// Out of range clamps rather than extrapolating.
	if got := Mix(Black, White, 2); got != White {
		t.Errorf("Mix(t=2) = %08X, want it clamped to the second colour", uint32(got))
	}
}

func TestFadeChangesOnlyAlpha(t *testing.T) {
	got := Fade(Red, 64)
	if got.A() != 64 {
		t.Errorf("alpha = %d, want 64", got.A())
	}
	if got&0x00FFFFFF != Red&0x00FFFFFF {
		t.Errorf("Fade changed the colour: %08X, want RGB of %08X", uint32(got), uint32(Red))
	}
}

func TestShade(t *testing.T) {
	if Shade(Gray, 0.5).R() <= Gray.R() {
		t.Error("lightening should raise the channels")
	}
	if Shade(Gray, -0.5).R() >= Gray.R() {
		t.Error("darkening should lower the channels")
	}
}

func Test565RoundTrip(t *testing.T) {
	if got := To565(White); got != 0xFFFF {
		t.Errorf("To565(White) = %04X, want FFFF", got)
	}
	// White must survive the round trip exactly: truncating instead of
	// repeating the top bits would bring it back at 248 and darken every
	// surface that is blended over itself.
	if got := From565(0xFFFF); got != White {
		t.Errorf("From565(FFFF) = %08X, want FFFFFFFF", uint32(got))
	}
	if got := From565(To565(Black)); got != Black {
		t.Errorf("black round trip = %08X, want FF000000", uint32(got))
	}
}
