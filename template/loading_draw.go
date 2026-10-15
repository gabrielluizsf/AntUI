package template

import (
	"time"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// draw is one frame of the screen, run from Window.End just before the
// pixels go out. When the time is up the app closes here — and the frame
// still draws, the bar at full, so the last thing on the display is a
// loading screen that reached its end rather than one cut in half. Ready
// may have been called while this very frame was being painted, and then
// there is nothing to draw: the app is ready, and the frame below is the
// program itself.
func (l *loading) draw(win *antui.Window, painter loadingPainter) {
	if win == nil || !l.active() {
		return
	}
	elapsed := time.Since(l.start)
	timeout := time.Duration(l.timeout.Load())
	if timeout > 0 && elapsed >= timeout {
		l.timedOut.Store(true)
		win.Quit()
	}
	win.Clear(painter.Background(win))
	// The bar sits in the middle of the window at three fifths of its width
	// and six units tall — a size that follows the window through
	// [Scale] rather than a number of pixels the screen picked.
	w := win.Width() * 3 / 5
	h := 6 * Scale(win)
	painter.Progress(win, canvas.Area{X: (win.Width() - w) / 2, Y: (win.Height() - h) / 2, Width: w, Height: h},
		progressAt(elapsed, timeout))
}
