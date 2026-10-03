package template

import (
	"sync/atomic"
	"time"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// defaultLoadingTimeout is how long a program may take to load before it
// closes itself, which is the number [SetLoadingTimeout] starts at and the
// one a program that says nothing gets: long enough for a stylesheet, a
// font and a database of icons, short enough that a program which never
// says Ready does not sit on a loading picture for ever.
const defaultLoadingTimeout = 5 * time.Second

// loadingEnabledByDefault says whether a template starts with its loading
// screen on. It is on for the program the developer ships — a screen the
// app has to ask to be without is the screen that catches the app which
// forgot to say Ready — and the test suite turns it off in an init, because
// the suite draws into canvases it then looks at and a screen painted over
// everything would be in the way of every one of those pictures. The tests
// that are about the loading screen turn it back on around their own cases.
var loadingEnabledByDefault = true

// loadingPainter is the half of a template the loading screen needs: the
// colour the window is cleared to, and the progress bar it draws with. Both
// [Template] and *uiTemplate answer to it, so the screen is drawn in the
// voice of whichever template the program picked — the same bar the rest of
// the app wears.
type loadingPainter interface {
	Background(win *antui.Window) canvas.Color
	Progress(win *antui.Window, x, y, w, h int, progress float32)
}

// loading is the start-up screen's state: when the wait began, how long it
// may last, and whether it is still up. Every field but the start time is
// atomic because the screen lives across goroutines — Ready arrives from
// whatever goroutine finished loading while the frames keep painting on the
// one that draws, and the two must not tear.
type loading struct {
	// start is when the template was made, which is when the screen first
	// appeared and the moment the timeout counts from. It is written once,
	// before the template is handed to anyone, and only read afterwards.
	start time.Time

	// timeout is the quit rule and the bar's denominator in nanoseconds, so
	// that both count against exactly the same number: the bar reaches 100%
	// at the instant the app would close. Zero or less is never stored —
	// [loading.setTimeout] leaves such a duration alone.
	timeout atomic.Int64

	// ready says the program called Ready, which ends the screen for good:
	// loading happens once, so nothing brings the screen back afterwards.
	ready atomic.Bool

	// disabled says the developer asked to be without the screen at all.
	// It starts as the negation of [loadingEnabledByDefault] when the
	// template is made.
	disabled atomic.Bool

	// pushed says the screen has been drawn and presented as the template's
	// first frame, which is what turning it off afterwards has to wipe.
	pushed atomic.Bool

	// timedOut says the timeout is what closed the app, kept so a test can
	// tell a program that ran out of time from one whose window was closed
	// some other way.
	timedOut atomic.Bool
}

// newLoading starts a screen at the default: on, counting from now, with
// the default timeout to close a program that never says Ready.
func newLoading() *loading {
	l := &loading{start: time.Now()}
	l.timeout.Store(int64(defaultLoadingTimeout))
	l.disabled.Store(!loadingEnabledByDefault)
	return l
}

// active says whether the screen is up: switched on and not finished. A nil
// screen — a template built as a bare struct rather than through [New] — is
// never active, which keeps every caller below free of a nil check.
func (l *loading) active() bool {
	return l != nil && !l.disabled.Load() && !l.ready.Load()
}

// markReady ends the screen. It may be called from any goroutine, which is
// the point: the goroutine that finished loading says so, and the next frame
// notices.
func (l *loading) markReady() {
	if l != nil {
		l.ready.Store(true)
	}
}

// setDisabled switches the screen off or back on. Coming back on does not
// rewind the clock: the timeout still counts from when the template was
// made, because that is when the screen first appeared.
func (l *loading) setDisabled(disabled bool) {
	if l != nil {
		l.disabled.Store(disabled)
	}
}

// setTimeout takes a new timeout, which is both how long the app may take
// and the span the bar fills over. A duration of zero or less would divide
// the bar by nothing and make the quit rule meaningless, so it is left out
// — an argument this package cannot use keeps the one already there, the
// way an unreadable attribute keeps the value it found.
func (l *loading) setTimeout(d time.Duration) {
	if l != nil && d > 0 {
		l.timeout.Store(int64(d))
	}
}

// getTimeout is the timeout in force.
func (l *loading) getTimeout() time.Duration {
	if l == nil {
		return 0
	}
	return time.Duration(l.timeout.Load())
}

// progress is how far along the bar is: the time gone over the time allowed,
// which puts 100% exactly where the timeout lands. It is read from the
// clock rather than ticked, so a frame that took long to draw still reports
// the truth and the bar never lies about how much time is left.
func (l *loading) progress() float32 {
	if l == nil {
		return 0
	}
	return progressAt(time.Since(l.start), time.Duration(l.timeout.Load()))
}

// progressAt is elapsed over timeout clamped into the range a bar reads.
// A timeout the API cannot produce (zero or less) has no denominator, and
// the bar stays where it is rather than inventing a number.
func progressAt(elapsed, timeout time.Duration) float32 {
	if timeout <= 0 {
		return 0
	}
	v := float64(elapsed) / float64(timeout)
	switch {
	case v < 0:
		return 0
	case v > 1:
		return 1
	default:
		return float32(v)
	}
}

// arm puts the screen on the frame being painted, by registering it to draw
// at the very end — [antui.Window.WidgetFrontmost] — above whatever the
// screen underneath asked for. Arm is called at the START of the frame's
// drawing, and Window.End runs the registered draws in reverse, so the
// loading screen lands last of all: above a dropdown, above a popup, above
// anything the app drew while it was still loading.
func (l *loading) arm(win *antui.Window, painter loadingPainter) {
	if win == nil || !l.active() {
		return
	}
	win.WidgetFrontmost(func() { l.draw(win, painter) })
}

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
	painter.Progress(win, (win.Width()-w)/2, (win.Height()-h)/2, w, h,
		progressAt(elapsed, timeout))
}

// showFirstFrame paints the screen the moment the template is made and puts
// it on the display, before the program has loaded anything. The window
// shows nothing until the first frame of the loop is presented, and
// everything a program does before that loop — reading a stylesheet, opening
// a database, fetching a picture — is exactly the stretch of time the screen
// is for, so the screen has to come before the loop rather than in it.
func (l *loading) showFirstFrame(win *antui.Window, painter loadingPainter) {
	if win == nil || !l.active() || l.pushed.Load() {
		return
	}
	l.pushed.Store(true)
	l.draw(win, painter)
	win.End()
}

// wipe puts the plain background where the screen was, and presents it. It
// is for the developer who turns the screen off after it already appeared —
// the flag set once the template exists — so that a program which asked to
// be without the loading picture does not sit on one until its first frame.
func (l *loading) wipe(win *antui.Window, painter loadingPainter) {
	if win == nil || l == nil || !l.pushed.Load() || l.ready.Load() {
		return
	}
	win.Clear(painter.Background(win))
	win.End()
}
