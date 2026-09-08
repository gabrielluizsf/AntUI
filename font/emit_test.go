package font

import (
	"testing"

	"github.com/gabrielluizsf/antui/assert"
)

func TestEmit_TooFewPoints(t *testing.T) {
	r := newRaster(10, 10)

	// Contour with only one point
	c := Contour{
		X:  []float64{5},
		Y:  []float64{5},
		On: []bool{true},
	}

	Emit(r, c, 1.0, 0, 0)
	mask := r.mask()

	// The mask should be completely empty (all 0s)
	hasFill := false
	for _, pixel := range mask {
		if pixel > 0 {
			hasFill = true
			break
		}
	}
	assert.False(t, hasFill)
}

func TestEmit_AllOnCurve(t *testing.T) {
	r := newRaster(10, 10)
	
	// A simple square shape using only on-curve points.
	c := Contour{
		X:  []float64{2, 8, 8, 2},
		Y:  []float64{2, 2, 8, 8},
		On: []bool{true, true, true, true},
	}

	// We pass oy=10.0 so the Y coordinates (2, 2, 8, 8) are mapped
	// to 8, 8, 2, 2 (oy - y), which fits inside our 10x10 raster.
	Emit(r, c, 1.0, 0, 10.0)
	mask := r.mask()
	
	hasFill := false
	for _, pixel := range mask {
		if pixel > 0 {
			hasFill = true
			break
		}
	}
	assert.True(t, hasFill)
}

func TestEmit_WithOffCurve(t *testing.T) {
	r := newRaster(10, 10)
	
	c := Contour{
		X:  []float64{2, 5, 8},
		Y:  []float64{2, 8, 2},
		On: []bool{true, false, true},
	}

	// Passing oy=10.0 to ensure coordinates fall inside the raster bounds.
	Emit(r, c, 1.0, 0, 10.0)
	mask := r.mask()
	
	hasFill := false
	for _, pixel := range mask {
		if pixel > 0 {
			hasFill = true
			break
		}
	}
	assert.True(t, hasFill)
}

func TestEmit_AllOffCurve(t *testing.T) {
	r := newRaster(20, 20)
	
	c := Contour{
		X:  []float64{5, 15, 15, 5},
		Y:  []float64{5, 5, 15, 15},
		On: []bool{false, false, false, false},
	}

	// Passing oy=20.0 to match the raster height.
	Emit(r, c, 1.0, 0, 20.0)
	mask := r.mask()
	
	hasFill := false
	for _, pixel := range mask {
		if pixel > 0 {
			hasFill = true
			break
		}
	}
	assert.True(t, hasFill)
}
func TestEmit_ScalingAndOffset(t *testing.T) {
	r := newRaster(10, 10)

	// A tiny square around the origin, which will be scaled and moved
	c := Contour{
		X:  []float64{0, 1, 1, 0},
		Y:  []float64{0, 0, 1, 1},
		On: []bool{true, true, true, true},
	}

	// scale = 5.0, ox = -2.0, oy = 8.0
	// This moves the tiny square into the visible 10x10 raster space
	Emit(r, c, 5.0, -2.0, 8.0)
	mask := r.mask()

	hasFill := false
	for _, pixel := range mask {
		if pixel > 0 {
			hasFill = true
			break
		}
	}
	assert.True(t, hasFill)
}
