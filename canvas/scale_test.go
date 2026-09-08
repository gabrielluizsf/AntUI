package canvas

import "testing"

// A picture reduced by an exact factor is the average of each block, and the
// averages here are chosen so that a wrong one is obvious rather than close.
func TestScaledAverages(t *testing.T) {
	src, err := NewCanvas(4, 4)
	if err != nil {
		t.Fatal(err)
	}
	// Four quadrants, each one solid colour.
	for y := range 4 {
		for x := range 4 {
			switch {
			case x < 2 && y < 2:
				src.Put(x, y, Black)
			case x >= 2 && y < 2:
				src.Put(x, y, White)
			case x < 2 && y >= 2:
				src.Put(x, y, RGB(0, 0, 255))
			default:
				src.Put(x, y, RGB(255, 0, 0))
			}
		}
	}
	got := Scaled(src, 2, 2)
	if got == nil {
		t.Fatal("Scaled returned nothing")
	}
	for _, want := range []struct {
		x, y int
		c    Color
	}{{0, 0, Black}, {1, 0, White}, {0, 1, RGB(0, 0, 255)}, {1, 1, RGB(255, 0, 0)}} {
		if c := got.At(want.x, want.y); c != want.c {
			t.Errorf("at (%d,%d) got %v, want %v", want.x, want.y, c, want.c)
		}
	}
}

// Scaling up repeats rather than inventing, so every pixel of the result is
// one that was in the original.
func TestScaledUpRepeats(t *testing.T) {
	src, _ := NewCanvas(2, 2)
	src.Put(0, 0, Black)
	src.Put(1, 0, White)
	src.Put(0, 1, White)
	src.Put(1, 1, Black)

	got := Scaled(src, 4, 4)
	for y := range 4 {
		for x := range 4 {
			want := src.At(x/2, y/2)
			if c := got.At(x, y); c != want {
				t.Fatalf("at (%d,%d) got %v, want %v", x, y, c, want)
			}
		}
	}
}

// IconScaled is Scaled to a square and must stay identical to it, or the two
// drift apart and an icon stops matching the picture it came from.
func TestIconScaledIsScaled(t *testing.T) {
	src, _ := NewCanvas(9, 5)
	for y := range 5 {
		for x := range 9 {
			src.Put(x, y, RGB(uint8(x*20), uint8(y*40), 128))
		}
	}
	icon, square := IconScaled(src, 4), Scaled(src, 4, 4)
	for y := range 4 {
		for x := range 4 {
			if icon.At(x, y) != square.At(x, y) {
				t.Fatalf("at (%d,%d): icon %v, scaled %v",
					x, y, icon.At(x, y), square.At(x, y))
			}
		}
	}
}

func TestFit(t *testing.T) {
	for _, tc := range []struct {
		srcW, srcH, boxW, boxH int
		wantW, wantH           int
	}{
		// A 16:9 picture in a tall box: the width fills and the height does
		// not, which is the case a phone screen always is.
		{1280, 720, 1080, 1920, 1080, 607},
		// The same picture in a wide box: the height fills instead.
		{1280, 720, 1920, 500, 888, 500},
		// Already the right shape, and already the right size.
		{100, 100, 50, 50, 50, 50},
		{100, 50, 100, 50, 100, 50},
		// Nothing sensible to answer.
		{0, 10, 10, 10, 0, 0},
	} {
		w, h := Fit(tc.srcW, tc.srcH, tc.boxW, tc.boxH)
		if w != tc.wantW || h != tc.wantH {
			t.Errorf("Fit(%d,%d in %d,%d) = %d,%d, want %d,%d",
				tc.srcW, tc.srcH, tc.boxW, tc.boxH, w, h, tc.wantW, tc.wantH)
		}
		if w > tc.boxW || h > tc.boxH {
			t.Errorf("Fit(%d,%d in %d,%d) = %d,%d, which does not fit",
				tc.srcW, tc.srcH, tc.boxW, tc.boxH, w, h)
		}
	}
}
