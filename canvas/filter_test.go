package canvas

import "testing"

func TestFilterGrayscaleAndInvert(t *testing.T) {
	cv, err := NewCanvas(4, 4)
	if err != nil {
		t.Fatal(err)
	}
	cv.FillRect(Area{0, 0, 4, 4}, Red)
	cv.FilterRegion(Area{X: 0, Y: 0, Width: 4, Height: 4}, FilterGrayscale, 1)
	g := cv.At(2, 2)
	if g.R() != g.G() || g.G() != g.B() {
		t.Fatalf("grayscale left colour: %v", g)
	}
	cv.FilterRegion(Area{X: 0, Y: 0, Width: 4, Height: 4}, FilterInvert, 1)
	if got := cv.At(2, 2); got.R() != 255-g.R() || got.A() != 255 {
		t.Fatalf("invert: %v from %v", got, g)
	}
}

func TestFilterBrightnessContrast(t *testing.T) {
	cv, _ := NewCanvas(2, 2)
	cv.FillRect(Area{0, 0, 2, 2}, RGBA(100, 100, 100, 255))
	cv.FilterRegion(Area{X: 0, Y: 0, Width: 2, Height: 2}, FilterBrightness, 0)
	if got := cv.At(0, 0); got.R() != 0 {
		t.Fatalf("brightness 0: %v", got)
	}
	cv.FillRect(Area{0, 0, 2, 2}, RGBA(100, 100, 100, 255))
	cv.FilterRegion(Area{X: 0, Y: 0, Width: 2, Height: 2}, FilterContrast, 2)
	if got := cv.At(0, 0); got.R() != 73 {
		t.Fatalf("contrast 2: got %d want 73", got.R())
	}
}

func TestFilterHueRotateFullTurnIsIdentity(t *testing.T) {
	cv, _ := NewCanvas(2, 2)
	cv.FillRect(Area{0, 0, 2, 2}, RGB(200, 40, 90))
	before := cv.At(0, 0)
	cv.FilterRegion(Area{X: 0, Y: 0, Width: 2, Height: 2}, FilterHueRotate, 360)
	after := cv.At(0, 0)
	for _, d := range []int{int(before.R()) - int(after.R()), int(before.G()) - int(after.G()), int(before.B()) - int(after.B())} {
		if d > 1 || d < -1 {
			t.Fatalf("full hue turn changed %v to %v", before, after)
		}
	}
}

func TestBlurSpreadsAndDims(t *testing.T) {
	cv, _ := NewCanvas(9, 9)
	cv.Put(4, 4, White)
	cv.Blur(Area{X: 0, Y: 0, Width: 9, Height: 9}, 2)
	if a := cv.At(4, 4).A(); a >= 255 {
		t.Fatalf("centre did not dim: alpha %d", a)
	}
	if a := cv.At(3, 4).A(); a == 0 {
		t.Fatal("blur did not reach a neighbour")
	}
}

func TestFilterOnNarrowCanvasIsNoOp(t *testing.T) {
	cv, err := NewCanvasFormat(2, 2, RGB565)
	if err != nil {
		t.Fatal(err)
	}
	cv.FilterRegion(Area{X: 0, Y: 0, Width: 2, Height: 2}, FilterInvert, 1)
	if cv.Pixels != nil {
		t.Fatal("narrow canvas should not have ARGB pixels")
	}
}
