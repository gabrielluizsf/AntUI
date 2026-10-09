package canvas

import (
	"fmt"
	"testing"
)

// TestBoxBlurIsTheSamePictureBitForBit holds the rewritten pass against the one
// it replaced. A blur is a blur whether it walks a line four times or once, but
// only if the arithmetic comes out the same: shadows, backdrop filters and
// blur() all land on screen through here, and a shadow that shifts by a shade
// from one frame to the next is a shadow that shimmers.
func TestBoxBlurIsTheSamePictureBitForBit(t *testing.T) {
	for _, size := range [][2]int{{1, 1}, {4, 3}, {17, 9}, {64, 33}} {
		for _, r := range []int{1, 2, 5, 40} {
			for _, horizontal := range []bool{true, false} {
				w, h := size[0], size[1]
				src := testPicture(w * h * 4)
				got := boxBlur(src, make([]float64, len(src)), w, h, r, horizontal)
				want := oneChannelAtATimeBlur(src, w, h, r, horizontal)
				for i := range got {
					if got[i] != want[i] {
						t.Fatalf("%dx%d radius %d horizontal %v: sample %d is %v, want %v",
							w, h, r, horizontal, i, got[i], want[i])
					}
				}
			}
		}
	}
}

// oneChannelAtATimeBlur is a box blur as it was written before: four separate
// walks of the line, one per channel, every window clamped. The slowest way to
// ask for the same numbers, and the one the fast path has to agree with.
func oneChannelAtATimeBlur(src []float64, w, h, r int, horizontal bool) []float64 {
	dst := make([]float64, len(src))
	span := float64(2*r + 1)
	if horizontal {
		for row := range h {
			base := row * w * 4
			for ch := range 4 {
				var sum float64
				for k := -r; k <= r; k++ {
					sum += src[base+clampInt(k, 0, w-1)*4+ch]
				}
				for col := range w {
					dst[base+col*4+ch] = sum / span
					sum -= src[base+clampInt(col-r, 0, w-1)*4+ch]
					sum += src[base+clampInt(col+r+1, 0, w-1)*4+ch]
				}
			}
		}
		return dst
	}
	for col := range w {
		for ch := range 4 {
			var sum float64
			for k := -r; k <= r; k++ {
				sum += src[clampInt(k, 0, h-1)*w*4+col*4+ch]
			}
			for row := range h {
				dst[row*w*4+col*4+ch] = sum / span
				sum -= src[clampInt(row-r, 0, h-1)*w*4+col*4+ch]
				sum += src[clampInt(row+r+1, 0, h-1)*w*4+col*4+ch]
			}
		}
	}
	return dst
}

// testPicture is a picture with enough going on in it — edges, a gradient, a
// transparent corner — that a blur of it is not the same as a blur of zeroes.
func testPicture(n int) []float64 {
	buf := make([]float64, n)
	for i := range buf {
		x := i % 4
		buf[i] = float64((i*37)%251)/7 + float64(x)*0.25
	}
	if n >= 4 {
		buf[0], buf[1], buf[2], buf[3] = 1, 0, 0, 0
	}
	return buf
}

// TestBlurKeepsTheEdgesHonest is what the clamping in a pass is for: a blur
// that ran off the edge of its picture would darken the border, and a border
// that darkens is a visible frame of dark pixels around every blurred box.
func TestBlurKeepsTheEdgesHonest(t *testing.T) {
	cv, err := NewCanvas(24, 24)
	if err != nil {
		t.Fatalf("NewCanvas: %v", err)
	}
	cv.FillRect(0, 0, 24, 24, RGB(200, 200, 200))
	cv.Blur(4, 4, 16, 16, 6)
	for _, p := range [][2]int{{4, 4}, {19, 4}, {4, 19}, {19, 19}, {12, 12}} {
		if c := cv.At(p[0], p[1]); c != RGB(200, 200, 200) {
			t.Errorf("inside pixel %v blurred to %v, want it untouched", p, c)
		}
	}
}

// TestBlurXYSpreadsOnlyWhereItIsAsked takes a single pixel and blurs it with a
// radius on one axis and none on the other: what it leaves is a line, not a
// square, which is the whole of what a per-axis radius is for.
func TestBlurXYSpreadsOnlyWhereItIsAsked(t *testing.T) {
	cv, err := NewCanvas(21, 21)
	if err != nil {
		t.Fatalf("NewCanvas: %v", err)
	}
	const cx, cy = 10, 10
	cv.Put(cx, cy, RGB(255, 255, 255))
	cv.BlurXY(0, 0, 21, 21, 0, 2)

	if c := cv.At(cx, cy-3); c.A() == 0 {
		t.Error("three pixels up is clear, want the blur carried there and not across")
	}
	if c := cv.At(cx+3, cy); c.A() != 0 {
		t.Errorf("three pixels across is %v, want nothing: the radius across is zero", c)
	}
}

// frame of an ordinary app draws: a blurred copy of a box the size of a button,
// laid over the page under it. It is bandwidth, not arithmetic — six passes of
// four channels of floats over the shadow's own pixels — so the numbers here
// are the reason the blur keeps its buffers between frames and the layers come
// from a pool.
func BenchmarkBoxShadow(b *testing.B) {
	for _, size := range [][2]int{{180, 60}, {200, 300}} {
		b.Run(fmt.Sprintf("%dx%d", size[0], size[1]), func(b *testing.B) {
			cv, err := NewCanvas(520, 620)
			if err != nil {
				b.Fatalf("NewCanvas: %v", err)
			}
			cv.FillRect(0, 0, 520, 620, RGB(0x10, 0x12, 0x18))
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				cv.BoxShadow(20, 20, size[0], size[1], 12, 12, 0, 8, 18, 0, RGBA(0, 0, 0, 0x99), false)
			}
		})
	}
}
