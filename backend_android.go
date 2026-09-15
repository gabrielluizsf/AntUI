//go:build android

package antui

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/display"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// The Android backend. A phone has no window manager to ask for anything: the
// system decides the size, when there is a surface at all, and when the app
// may draw. So most of the backend interface is a polite refusal, and the
// two calls that do the work are pump — which waits when there is nothing to
// draw on — and present, which is a blit into a buffer the compositor lends
// out one frame at a time.
//
// The Canvas above this is the same Canvas as everywhere else, so a program
// written for a desktop draws on a phone without knowing it moved.
type androidWindow struct {
	activity *app.App
	surface  *ndk.Window
	format   ndk.Format
	scale    float64

	// primary is the touch that is also reported as the left mouse button,
	// or -1 when no finger is down. buttons is the last mouse button state
	// seen, which is what a change is worked out against.
	primary int
	buttons int32
	// refresh is the display's rate, asked once. Zero means not yet asked
	// and a negative means asked and not answered.
	refresh int
}

func newBackend() platform { return &androidWindow{primary: -1} }

// errNotAndroidApp is what an Android *process* gets when it is not an
// Android *app* — a command-line binary pushed to a device with adb, say,
// which has a Go runtime and no activity.
var errNotAndroidApp = errors.New("antui: no Android activity: this process was not started by the system as an app")

// open waits for the system to give the app a surface, and takes its size
// from it. The width and height asked for are ignored, because there is
// nothing to ask: a phone's window is the screen.
func (b *androidWindow) open(win *Window, title string, width, height int) error {
	b.activity = app.Current()
	if b.activity == nil {
		return errNotAndroidApp
	}
	b.scale = scaleFromDensity(b.activity.Density())
	// There is no other kind of input here, so a game does not have to wait
	// for a finger to know it needs controls on the screen.
	win.SetTouchFirst()

	// The activity exists before its surface does. Everything up to the
	// first WindowUp is handled here rather than being dropped, so that a
	// Start and a Resume that arrive first are not lost.
	for b.surface == nil {
		e, ok := b.activity.Next()
		if !ok {
			return errors.New("antui: the activity was destroyed before it had a surface")
		}
		b.handle(win, e)
		if win.shouldClose {
			return errors.New("antui: the activity closed before it had a surface")
		}
	}
	b.activity.Release()
	b.sync(win)
	return nil
}

// close asks Android to finish the activity. The process is not ended here:
// that is the system's to do, and a native app that calls exit is a native
// app that loses whatever the platform was about to save.
func (b *androidWindow) close() {
	b.surface = nil
	app.Finish()
}

// pump takes every event that is waiting, and then — if there is nothing to
// draw on — waits for one.
//
// That wait is the whole of what makes an app behave when it is not in
// front. A frame loop with no surface has nothing to do, and spinning
// through it burns a battery and gets the process killed; parked on a
// channel it costs nothing until the app is resumed.
func (b *androidWindow) pump(win *Window) {
	if b.activity == nil {
		return
	}
	for {
		e, ok := b.activity.Poll()
		if !ok {
			break
		}
		b.handle(win, e)
	}
	b.drainInput(win)
	for b.surface == nil && !win.shouldClose {
		e, ok := b.activity.Next()
		if !ok {
			win.Push(Event{Type: EventClose})
			break
		}
		b.handle(win, e)
	}
	// Let go of the last event before the frame starts. Everything after
	// this belongs to the app, and the UI thread must not be blocked for the
	// length of it.
	b.activity.Release()
	b.sync(win)
}

// handle turns one activity event into what the window core understands.
func (b *androidWindow) handle(win *Window, e app.Event) {
	switch e {
	case app.WindowUp:
		b.surface = b.activity.Window()
		if b.surface != nil {
			// Ask for eight bits a channel rather than take what the surface
			// offers: a device that defaults to 565 would otherwise quietly
			// halve the colour depth of everything drawn.
			if err := b.surface.SetGeometry(0, 0, ndk.RGBA8888); err != nil {
				ndk.Warnf("cannot set the surface format: %v", err)
			}
			b.format = b.surface.Format()
		}
		win.PushSimple(EventExpose)

	case app.WindowDown:
		// The window is already gone from the activity; dropping it here as
		// well is what makes present a no-op instead of a crash.
		b.surface = nil

	case app.Redraw, app.WindowResized:
		win.PushSimple(EventExpose)

	case app.FocusGained:
		win.Push(Event{Type: EventFocus, Focused: true})
	case app.FocusLost:
		win.Push(Event{Type: EventFocus, Focused: false})

	case app.ContentRect:
		b.safe(win)

	case app.ConfigChanged:
		b.scale = scaleFromDensity(b.activity.Density())
		b.safe(win)
		win.PushSimple(EventExpose)

	case app.Destroy:
		b.surface = nil
		win.Push(Event{Type: EventClose})
	}
}

// safe passes the system's content rectangle on as the window's safe area.
//
// It comes from the platform for nothing — onContentRectChanged is a
// callback, not a call — which is why it is used rather than the window
// insets, which are a Java round trip. The insets say more (they separate
// the keyboard and the cutout) and app.WindowInsets is there for a program
// that needs the difference.
func (b *androidWindow) safe(win *Window) {
	r := b.activity.ContentRect()
	if r.Right <= r.Left || r.Bottom <= r.Top {
		return
	}
	win.SetSafeArea(canvas.Area{
		X:      int(r.Left),
		Y:      int(r.Top),
		Width:  int(r.Right - r.Left),
		Height: int(r.Bottom - r.Top),
	})
}

// sync makes the canvas the size of the surface. On a desktop a resize is
// news from the window manager; here it is simply a fact to be noticed,
// because nothing was ever asked for.
func (b *androidWindow) sync(win *Window) {
	if b.surface == nil {
		return
	}
	w, h := b.surface.Width(), b.surface.Height()
	if w <= 0 || h <= 0 || (w == win.cv.Width && h == win.cv.Height) {
		return
	}
	if win.ResizeCanvas(w, h) {
		win.Push(Event{Type: EventResize, Width: w, Height: h})
		// resizeCanvas throws the safe area away, since it was in the old
		// canvas's coordinates. The platform's rectangle is still current.
		b.safe(win)
	}
}

// present blits the frame into the buffer the compositor lends us.
//
// The rectangle handed in is what changed. The compositor's answer to it is
// usually larger — it hands out buffers in rotation, so the one being locked
// holds a frame from two or three ago, and everything stale in it has to be
// written. Lock rewrites the rectangle with what must be filled, and that is
// the one to blit.
func (b *androidWindow) present(win *Window, dirty canvas.Area) {
	if b.surface == nil || dirty.Width <= 0 || dirty.Height <= 0 {
		return
	}
	r := ndk.Rect{
		Left:   int32(dirty.X),
		Top:    int32(dirty.Y),
		Right:  int32(dirty.X + dirty.Width),
		Bottom: int32(dirty.Y + dirty.Height),
	}
	buf, err := b.surface.Lock(&r)
	if err != nil {
		// The surface can be taken away between the event and this call.
		// Losing a frame is the right answer; the WindowDown event is on its
		// way and pump will park.
		ndk.Warnf("cannot lock the surface: %v", err)
		b.surface = nil
		return
	}
	blit(buf, win.cv, r)
	if err := b.surface.UnlockAndPost(); err != nil {
		ndk.Warnf("cannot post the frame: %v", err)
		b.surface = nil
	}
}

// setFullscreen hides the status bar. The navigation bar stays: taking that
// away is immersive mode, which Android only offers through Java.
func (b *androidWindow) setFullscreen(on bool) bool {
	if on {
		app.SetWindowFlags(app.FlagFullscreen, 0)
	} else {
		app.SetWindowFlags(0, app.FlagFullscreen)
	}
	return true
}

// displaySize is the surface, which on a phone is the display.
func (b *androidWindow) displaySize() (int, int, bool) {
	if b.surface == nil {
		return 0, 0, false
	}
	w, h := b.surface.Width(), b.surface.Height()
	return w, h, w > 0 && h > 0
}

// displayRefresh is the panel's rate, asked once and remembered: it comes
// from Java, and a round trip to the machine for every caller of
// Window.DisplayRefresh would be a poor way to answer a number that changes
// about as often as the hardware does.
//
// It is rounded, because this is where the honest float has to become the
// integer antui's API promised. 59.94 is 60.
func (b *androidWindow) displayRefresh() int {
	if b.refresh == 0 {
		b.refresh = -1 // asked and unanswerable, so do not ask again
		if hz, err := display.Refresh(); err == nil && hz > 0 {
			b.refresh = int(hz + 0.5)
		}
	}
	if b.refresh < 0 {
		return 0
	}
	return b.refresh
}

// contentScale is the density over 160, which is the number Android itself
// multiplies a size in device-independent pixels by.
func (b *androidWindow) contentScale() float64 { return b.scale }

func scaleFromDensity(dpi int) float64 {
	if dpi <= 0 {
		return 0
	}
	return float64(dpi) / 160
}

// The rest cannot be done on a phone, and say so rather than pretending.

// setTitle: an Android app's name is in its manifest, and the label in the
// recents list is set through Java. Nothing here can change it.
func (b *androidWindow) setTitle(string) {}

// setLimits and setSize: the size of the window is the system's decision and
// is not negotiable. Reporting false is the honest answer, and window_size.go
// is written to accept it.
func (b *androidWindow) setLimits(*Window, Limits)      {}
func (b *androidWindow) setSize(*Window, int, int) bool { return false }

// setIcon: the icon is a resource in the package, chosen when the APK is
// built and not at runtime.
func (b *androidWindow) setIcon([]*canvas.Canvas) bool { return false }

// The clipboard is Java's, and arrives with the JNI layer.
func (b *androidWindow) clipboard() (string, []string) { return "", nil }
func (b *androidWindow) setClipboard(string) bool      { return false }
