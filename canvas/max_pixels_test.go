package canvas

import (
	"math"
	"testing"
)

// TestNewCanvasRefusesASizeThatWraps is a width and a height whose product
// passes the top of an int and comes back round to almost nothing: the canvas
// would say it was enormous while holding a few pixels, and the first row
// written into it would run off the end.
func TestNewCanvasRefusesASizeThatWraps(t *testing.T) {
	for _, size := range [][2]int{
		{1 << 32, 1 << 32},
		{math.MaxInt, 2},
		{MaxPixels, 2},
		{1 << 20, 1 << 20},
	} {
		if cv, err := NewCanvas(size[0], size[1]); err == nil {
			t.Errorf("NewCanvas(%d, %d) made a canvas of %d pixels, want it refused", size[0], size[1], len(cv.Pixels))
		}
		if _, err := NewCanvasFormat(size[0], size[1], RGB565); err == nil {
			t.Errorf("NewCanvasFormat(%d, %d) made a canvas, want it refused", size[0], size[1])
		}
		if _, ok := acquireLayer(size[0], size[1]); ok {
			t.Errorf("acquireLayer(%d, %d) made a layer, want it refused", size[0], size[1])
		}
	}
}

// TestNewCanvasStillMakesABigScreen keeps the limit above anything a program
// draws: an 8K screen is made as before.
func TestNewCanvasStillMakesABigScreen(t *testing.T) {
	if _, err := NewCanvas(7680, 4320); err != nil {
		t.Fatalf("an 8K canvas was refused: %v", err)
	}
}
