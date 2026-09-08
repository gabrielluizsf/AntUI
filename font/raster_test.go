package font

import (
	"testing"

	"github.com/gabrielluizsf/antui/assert"
)

func TestNewRaster(t *testing.T) {
	// Valid dimensions
	r := newRaster(10, 20)
	assert.Equal(t, 10, r.w)
	assert.Equal(t, 20, r.h)
	// (w+1)*h = 11 * 20 = 220
	assert.Len(t, 220, r.a)

	// Invalid dimensions (zero or negative)
	rZero := newRaster(0, 0)
	assert.Equal(t, 0, rZero.w)
	assert.Equal(t, 0, rZero.h)
	assert.Len(t, 0, rZero.a)

	rNegative := newRaster(-5, -5)
	assert.Equal(t, 0, rNegative.w)
	assert.Equal(t, 0, rNegative.h)
	assert.Len(t, 0, rNegative.a)
}

func TestRaster_Mask_Empty(t *testing.T) {
	r := newRaster(0, 0)
	mask := r.mask()
	assert.Len(t, 0, mask)
}

func TestRaster_LineAndMask(t *testing.T) {
	// Create a 4x4 image
	r := newRaster(4, 4)

	// Draw a 2x2 square in the middle of the 4x4 grid.
	// We draw the left edge going down, and the right edge going up.
	// Horizontal lines are ignored by the rasterizer (y0 == y1), 
	// so the vertical edges define the fill via winding.
	
	// Left edge (x=1), going down from y=1 to y=3
	r.line(1, 1, 1, 3)
	// Right edge (x=3), going up from y=3 to y=1
	r.line(3, 3, 3, 1)

	mask := r.mask()
	// 4 * 4 = 16 pixels
	assert.Len(t, 16, mask)

	// Check y=0 (Row 0) - Should be completely empty
	assert.Equal(t, uint8(0), mask[0])
	assert.Equal(t, uint8(0), mask[1])
	assert.Equal(t, uint8(0), mask[2])
	assert.Equal(t, uint8(0), mask[3])

	// Check y=1 (Row 1) - Pixels at x=1 and x=2 should be filled
	assert.Equal(t, uint8(0), mask[4])
	assert.Equal(t, uint8(255), mask[5])
	assert.Equal(t, uint8(255), mask[6])
	assert.Equal(t, uint8(0), mask[7])

	// Check y=2 (Row 2) - Pixels at x=1 and x=2 should be filled
	assert.Equal(t, uint8(0), mask[8])
	assert.Equal(t, uint8(255), mask[9])
	assert.Equal(t, uint8(255), mask[10])
	assert.Equal(t, uint8(0), mask[11])

	// Check y=3 (Row 3) - Should be completely empty
	assert.Equal(t, uint8(0), mask[12])
	assert.Equal(t, uint8(0), mask[13])
	assert.Equal(t, uint8(0), mask[14])
	assert.Equal(t, uint8(0), mask[15])
}

func TestRaster_Quad(t *testing.T) {
	r := newRaster(10, 10)

	// Draw a curve from (1,1) to (9,1) with a control point at (5,9)
	r.quad(1, 1, 5, 9, 9, 1)
	
	// Close the shape with a straight line back to the start to form a solid
	r.line(9, 1, 1, 1)

	mask := r.mask()
	assert.Len(t, 100, mask)

	// Check if any pixels were drawn (coverage > 0)
	hasFill := false
	for _, pixel := range mask {
		if pixel > 0 {
			hasFill = true
			break
		}
	}
	
	assert.True(t, hasFill)
}

func TestRaster_Span(t *testing.T) {
	r := newRaster(5, 5)
	
	// y=2, xa=1.5, xb=3.5, height=1.0
	// It should cover part of pixel 1, all of pixel 2, and part of pixel 3.
	r.span(2, 1.5, 3.5, 1.0)
	
	mask := r.mask()
	assert.Len(t, 25, mask)
	
	// Pixel x=0 at y=2
	assert.Equal(t, uint8(0), mask[2*5+0])
	
	// Due to accumulated coverage spanning, there should be positive values 
	// from x=1 to x=3 on y=2.
	assert.True(t, mask[2*5+1] > 0)
	assert.True(t, mask[2*5+2] > 0)
	assert.True(t, mask[2*5+3] > 0)
}