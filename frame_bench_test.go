package antui

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// BenchmarkFrameRedrawn measures a frame that paints the same picture it
// painted last: the window sitting still, the clock not moving, nothing
// hovered. Every pixel is written and every written row is measured against
// the frame on the screen, and the answer is that nothing changed — so no
// pixels are converted, no request is built and no socket is touched.
//
// The number to watch is the allocations: a steady frame that hands the
// collector a frame's worth of garbage is a frame that stutters every time the
// heap grows.
func BenchmarkFrameRedrawn(b *testing.B) {
	win, _ := newTestWindow(b, 1280, 800)
	win.SetFPS(0)
	win.sleep = func(time.Duration) {}

	// The first frame is the whole window: there is nothing on the screen yet.
	drawTestFrame(win)
	win.dirtyRegion()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		win.Begin()
		drawTestFrame(win)
		if win.End() {
			b.Fatal("a frame that redraws the same picture presented something")
		}
	}
}

// BenchmarkFrameWithACursor measures the frame a program draws when one pixel
// of it moved: a mouse over a still window. The damage is six pixels, and the
// measuring is bounded by the rows that were written rather than by the size of
// the canvas, so this stays close to the cost of drawing.
func BenchmarkFrameWithACursor(b *testing.B) {
	win, _ := newTestWindow(b, 1280, 800)
	win.SetFPS(0)
	win.sleep = func(time.Duration) {}

	drawTestFrame(win)
	win.dirtyRegion()

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		win.Begin()
		drawTestFrame(win)
		win.Clear(canvas.RGB(0x10, 0x12, 0x18))
		win.FillRect(100+i%900, 100, 4, 4, canvas.White)
		if !win.End() {
			b.Fatal("a frame with a cursor on it presented nothing")
		}
	}
}

// BenchmarkFrameAnimated measures a whole window repainted every frame with
// something moving in it, which is the honest cost of a window that animates:
// everything is written, and the whole window goes to the display.
func BenchmarkFrameAnimated(b *testing.B) {
	win, _ := newTestWindow(b, 1280, 800)
	win.SetFPS(0)
	win.sleep = func(time.Duration) {}

	drawTestFrame(win)
	win.dirtyRegion()

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		win.Begin()
		drawTestFrame(win)
		win.FillRect(0, i%700, 1280, 8, canvas.RGB(0xFF, 0x88, 0x00))
		win.End()
	}
}

// BenchmarkPace measures the loop's own timing: the frames land on a grid, so
// a frame that took too long is followed by a short wait rather than a long
// one, and a frame that overran drops the next instead of accumulating a debt.
func BenchmarkPace(b *testing.B) {
	win, _ := newTestWindow(b, 64, 64)
	win.SetFPS(60)
	var waited time.Duration
	win.sleep = func(d time.Duration) { waited += d }

	b.ResetTimer()
	for range b.N {
		waited = 0
		win.pace()
	}
	b.StopTimer()
	if waited <= 0 {
		b.Fatal("the loop never waited")
	}
}

// drawTestFrame paints a window's worth of ordinary interface: a background, a
// few blocks, some text and a ring. It is the same every frame on purpose, so
// the benchmarks above measure what happens when nothing moves.
func drawTestFrame(win *Window) {
	win.Clear(canvas.RGB(0x10, 0x12, 0x18))
	win.FillRect(20, 20, 300, 180, canvas.RGB(0x1E, 0x24, 0x30))
	win.FillRoundRect(340, 20, 180, 60, 10, canvas.RGB(0x3E, 0x63, 0xDD))
	win.RoundRect(20, 220, 300, 40, 8, canvas.RGB(0x6B, 0x74, 0x88))
	win.Text(40, 50, "AntUI builds desktop apps the simple way.", canvas.White)
	win.Circle(700, 300, 120, canvas.RGB(0x23, 0x2A, 0x38))
	win.Canvas().RingArc(700, 300, 120, 14, 0, 2.0, canvas.RGB(0x3E, 0x63, 0xDD))
	win.Line(0, 799, 1279, 0, canvas.RGBA(0xFF, 0xFF, 0xFF, 0x40))
}
