package canvas

import (
	"math"
	"testing"
)

func off(c Color) GradientStop { return GradientStop{Offset: -1, Color: c} }

func TestLinearGradientRunsTopToBottom(t *testing.T) {
	cv, err := NewCanvas(20, 20)
	if err != nil {
		t.Fatal(err)
	}
	// to bottom: the first stop lands on the top side, the last on the bottom.
	cv.FillGradient(0, 0, 20, 20, Gradient{
		Kind:  GLinear,
		Angle: math.Pi,
		Stops: []GradientStop{off(RGBA(255, 0, 0, 255)), off(RGBA(0, 0, 255, 255))},
	})
	if got := cv.At(10, 0); got.A() == 0 || got == cv.At(10, 19) {
		t.Errorf("top should be the first stop red, got %v", got)
	}
	if got := cv.At(10, 19); got.B() < got.R() {
		t.Errorf("bottom pixel = %v, want it blue-dominant", got)
	}
}

func TestLinearGradientRunsLeftToRight(t *testing.T) {
	cv, _ := NewCanvas(20, 20)
	cv.FillGradient(0, 0, 20, 20, Gradient{
		Kind:  GLinear,
		Angle: math.Pi / 2, // to right
		Stops: []GradientStop{off(RGBA(255, 0, 0, 255)), off(RGBA(0, 0, 255, 255))},
	})
	if got := cv.At(0, 10); got.A() == 0 || got == cv.At(19, 10) {
		t.Errorf("left should be the first stop red, got %v", got)
	}
	if got := cv.At(19, 10); got.B() < got.R() {
		t.Errorf("right pixel = %v, want it blue-dominant", got)
	}
}

func TestGradientSpreadsUnsetStopsEvenly(t *testing.T) {
	cv, _ := NewCanvas(40, 40)
	cv.FillGradient(0, 0, 40, 40, Gradient{
		Kind: GLinear,
		Stops: []GradientStop{
			off(RGBA(255, 0, 0, 255)),
			off(RGBA(0, 255, 0, 255)),
			off(RGBA(0, 0, 255, 255)),
		},
	})
	mid := cv.At(10, 20)
	redRow := cv.At(10, 10)
	if mid == redRow || mid.A() == 0 {
		t.Errorf("the midpoint of the three stops should blend toward the middle colour, got %v", mid)
	}
}

func TestSingleStopGradientIsSolid(t *testing.T) {
	cv, _ := NewCanvas(20, 20)
	cv.FillGradient(0, 0, 20, 20, Gradient{
		Kind:  GLinear,
		Stops: []GradientStop{off(RGBA(1, 2, 3, 255))},
	})
	if got, want := cv.At(10, 10), RGBA(1, 2, 3, 255); got != want {
		t.Errorf("single stop should paint solid, got %v want %v", got, want)
	}
}

func TestRadialGradientCentred(t *testing.T) {
	cv, _ := NewCanvas(21, 21)
	cv.FillGradient(0, 0, 21, 21, Gradient{
		Kind:    GRadial,
		CenterX: 0.5,
		CenterY: 0.5,
		Stops:   []GradientStop{off(RGBA(255, 0, 0, 255)), off(RGBA(0, 0, 255, 255))},
	})
	cx := cv.At(10, 10)
	e := cv.At(20, 10)
	if cx.A() == 0 || e.A() == 0 {
		t.Fatalf("radial gradient left holes: centre=%v edge=%v", cx, e)
	}
	if cx == e {
		t.Errorf("radial centre and far edge should differ, both %v", cx)
	}
}

func TestConicStartsPointingUp(t *testing.T) {
	cv, _ := NewCanvas(21, 21)
	cv.FillGradient(0, 0, 21, 21, Gradient{
		Kind:    GConic,
		CenterX: 0.5,
		CenterY: 0.5,
		Stops:   []GradientStop{off(RGBA(255, 0, 0, 255)), off(RGBA(0, 0, 255, 255))},
	})
	top := cv.At(10, 1)     // overhead: the start angle, first stop
	bottom := cv.At(10, 20) // a half turn back round
	if top == bottom {
		t.Errorf("conic opposite sides should be the two stops, got %v both", top)
	}
}

func TestMaskRoundRectKeepsInsideZerosOutside(t *testing.T) {
	cv, _ := NewCanvas(30, 30)
	cv.FillRect(0, 0, 30, 30, White)
	cv.MaskRoundRect(5, 5, 20, 20, 6, 6)
	if got := cv.At(2, 2).A(); got != 0 {
		t.Errorf("outside the mask should be cleared, got alpha %d", got)
	}
	if got := cv.At(12, 12).A(); got != 255 {
		t.Errorf("inside the mask should survive, got alpha %d", got)
	}
	if got := cv.At(10, 2).A(); got != 0 {
		t.Errorf("above the mask should be cleared, got alpha %d", got)
	}
}

func TestMaskRectKeepsInsideZerosOutside(t *testing.T) {
	cv, _ := NewCanvas(10, 10)
	cv.FillRect(0, 0, 10, 10, White)
	cv.MaskRect(3, 3, 4, 4)
	if got := cv.At(1, 1).A(); got != 0 {
		t.Errorf("outside should be cleared, got alpha %d", got)
	}
	if got := cv.At(4, 4).A(); got != 255 {
		t.Errorf("inside should survive, got alpha %d", got)
	}
}
