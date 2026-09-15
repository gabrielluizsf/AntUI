//go:build android

package app

import (
	"sync"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Event is one thing that happened to the app. They arrive in the order
// Android raised them, one at a time, from [App.Next].
type Event int

// The events. The ones that carry something say so.
const (
	// None is no event: what [App.Poll] returns when nothing is waiting.
	None Event = iota

	// Start, Resume, Pause and Stop are the activity lifecycle. An app is
	// between Resume and Pause when the user is looking at it.
	Start
	Resume
	Pause
	Stop
	// Destroy is the last event. Next reports false after it.
	Destroy

	// WindowUp says there is a surface to draw on: [App.Window] is not nil
	// from here until WindowDown.
	WindowUp
	// WindowResized and Redraw both mean draw a whole frame now. Redraw in
	// particular is the system saying the old contents are gone.
	WindowResized
	Redraw
	// WindowDown says the surface is going away. **The UI thread is blocked
	// until the next call to Next**, so stop drawing and come back promptly;
	// [App.Window] is already nil.
	WindowDown

	// FocusGained and FocusLost are the window's keyboard focus.
	FocusGained
	FocusLost

	// ContentRect is the area not covered by system bars; [App.ContentRect]
	// has it.
	ContentRect
	// ConfigChanged is a rotation, a locale, a density, a theme.
	ConfigChanged
	// LowMemory is the system asking for caches back.
	LowMemory

	// InputUp and InputDown bracket the input queue's life. Reading it is
	// the input layer's business, not this package's.
	InputUp
	InputDown
)

var eventNames = [...]string{
	"None",
	"Start", "Resume", "Pause", "Stop", "Destroy",
	"WindowUp", "WindowResized", "Redraw", "WindowDown",
	"FocusGained", "FocusLost",
	"ContentRect", "ConfigChanged", "LowMemory",
	"InputUp", "InputDown",
}

// String names the event, which is what a log line wants.
func (e Event) String() string {
	if int(e) < len(eventNames) {
		return eventNames[e]
	}
	return "Event(?)"
}

// App is the running activity. There is only ever one, and [Main] is given
// it when the system creates it.
type App struct {
	msgs chan message
	// dead is closed when the app goroutine will never receive again, so a
	// callback on the UI thread gives up instead of blocking forever.
	dead chan struct{}

	// Everything below is written on the UI thread and read on the app's,
	// so it is behind the mutex even when it looks like it could not be.
	mu     sync.Mutex
	window *ndk.Window
	rect   ndk.Rect
	focus  bool
	info   Info
	config Config
	state  State
	// saved is what a previous run of this activity wrote down; see
	// state.go.
	saved []byte

	// Everything below belongs to the app's goroutine alone and is never
	// touched from a callback, so none of it is behind the mutex.
	pending message
	held    bool
	done    bool
	// loop is this thread's looper, and queue the input attached to it.
	// Both are thread-local in the platform's sense, which is why the
	// goroutine is pinned to an OS thread.
	loop  *ndk.Looper
	queue *ndk.InputQueue

	// assets is the manager the platform handed the activity. It is set
	// once, before the app's goroutine starts, and never changes.
	assets unsafe.Pointer
}

// Assets is the reader for the files packed into the APK.
func (a *App) Assets() *ndk.AssetManager {
	if a == nil {
		return nil
	}
	return ndk.WrapAssets(a.assets)
}

// Info is what the activity knows about itself and never changes.
type Info struct {
	// InternalData is the app's private directory: what os.UserCacheDir
	// would be if there were a home directory, and the only place writing is
	// certain to work.
	InternalData string
	// ExternalData is the app's directory on shared storage, which may not
	// exist and may not be writable.
	ExternalData string
	// OBB is where an expansion file would be mounted.
	OBB string
	// Cache is where things that can be thrown away go. The system deletes
	// it when the device runs short, without asking and without warning, so
	// nothing that matters belongs here — and everything that does not
	// belongs nowhere else.
	Cache string
	// SDK is the API level of the device — 21 and up.
	SDK int
}

type message struct {
	ev  Event
	ptr unsafe.Pointer
	// ack is closed by the app goroutine when the event has been dealt with.
	// It is nil for the events the UI thread does not wait on.
	ack chan struct{}
}

// Next is the next event, and false once the activity is destroyed.
//
// It also acknowledges the event before it. Android blocks its UI thread on
// a few of these — WindowDown above all — and the block is released here, on
// the *following* call, or by [App.Release]. That is the contract: an event
// is in hand from the Next that returned it until the next Next, Poll or
// Release, and nothing that the event took away may be touched during it.
func (a *App) Next() (Event, bool) {
	if a.done {
		return None, false
	}
	a.release()
	m, ok := <-a.msgs
	if !ok {
		a.done = true
		return None, false
	}
	return a.take(m), true
}

// Poll is Next for a caller that cannot afford to wait: it returns the next
// event if one is already there, and None with false if none is. The
// acknowledgement rule is the same — the event before is released here, so a
// frame loop that polls once a frame holds the UI thread for at most a frame.
//
// A false result means "nothing right now", not "the app is over".
// [App.Alive] is what says that.
func (a *App) Poll() (Event, bool) {
	if a.done {
		return None, false
	}
	a.release()
	select {
	case m, ok := <-a.msgs:
		if !ok {
			a.done = true
			return None, false
		}
		return a.take(m), true
	default:
		return None, false
	}
}

// take is the half of Next and Poll that is the same: hold the event, fold
// it into the state, and notice the last one.
func (a *App) take(m message) Event {
	a.pending, a.held = m, true
	a.apply(m)
	a.attach(m)
	if m.ev == Destroy {
		a.done = true
		a.release()
		close(a.dead)
	}
	return m.ev
}

// attach hooks the input queue up to this thread's looper, and unhooks it.
//
// It has to happen here and not in a callback: a queue is delivered to the
// looper of the thread that attached it, and the UI thread's looper is not
// the one this app reads. Attaching from the wrong thread does not fail — it
// simply means no input ever arrives.
func (a *App) attach(m message) {
	switch m.ev {
	case InputUp:
		if a.loop == nil {
			a.loop = ndk.PrepareLooper()
		}
		a.queue = ndk.WrapQueue(m.ptr)
		if a.loop != nil && a.queue != nil {
			a.queue.AttachLooper(a.loop)
		}
	case InputDown:
		if a.queue != nil {
			a.queue.DetachLooper()
			a.queue = nil
		}
	}
}

// Input takes everything waiting in the input queue and hands each event to
// handle, which reports whether the app dealt with it. An event the app does
// not deal with goes back to the system — which is how the back button still
// closes an app that ignored it.
//
// It must be called from the same goroutine as [App.Next], which is the one
// the looper belongs to.
func (a *App) Input(handle func(*ndk.InputEvent) bool) {
	if a.queue == nil || a.loop == nil {
		return
	}
	// Give the looper a turn so it notices the queue is readable. Emptying
	// it is the loop below; this only wakes the machinery.
	a.loop.Poll(0)
	for {
		e, ok := a.queue.Next()
		if !ok {
			return
		}
		// The system gets first refusal on everything. This is how the
		// on-screen keyboard sees a key at all, and skipping it produces an
		// app that works until somebody types.
		if a.queue.PreDispatch(e) {
			continue
		}
		a.queue.Finish(e, handle(e))
	}
}

// Alive reports whether the activity still exists. It goes false once
// Destroy has been taken, and stays false.
func (a *App) Alive() bool { return a != nil && !a.done }

// Release lets go of the event in hand, freeing the UI thread if it was
// waiting on that one.
//
// [App.Next] and [App.Poll] do this for the event before, so a loop that
// keeps asking never needs it. It matters for a loop that stops asking and
// goes off to do something else — a frame — because the UI thread would
// otherwise stay blocked for the whole of it. Blocked, in particular,
// against [RunOnUISync], which would then be waiting for a thread that is
// waiting for it.
func (a *App) Release() { a.release() }

// release lets the UI thread go, if it was waiting on the event in hand.
func (a *App) release() {
	if !a.held {
		return
	}
	a.held = false
	if a.pending.ack != nil {
		close(a.pending.ack)
		a.pending.ack = nil
	}
}

// apply folds an event into the app's state before the caller sees it — so
// that on WindowDown the window is already gone, and on WindowUp it is
// already there.
func (a *App) apply(m message) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.noteState(m.ev)
	switch m.ev {
	case WindowUp:
		a.window = ndk.Wrap(m.ptr)
	case WindowDown:
		a.window = nil
	case FocusGained:
		a.focus = true
	case FocusLost:
		a.focus = false
	}
}

// Window is the surface to draw on, or nil when there is none. It is only
// valid between WindowUp and WindowDown, and the whole point of Next's
// contract is that "between" means what it says.
func (a *App) Window() *ndk.Window {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.window
}

// ContentRect is the part of the window not hidden behind a status bar, a
// navigation bar or a cutout.
func (a *App) ContentRect() ndk.Rect {
	if a == nil {
		return ndk.Rect{}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.rect
}

// Focused reports whether the window has focus.
func (a *App) Focused() bool {
	if a == nil {
		return false
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.focus
}

// Info is what the activity knows about itself.
//
// A nil App answers with a zero Info rather than crashing. [Current] returns
// nil while there is no activity — before one exists, and for the moment
// after one is destroyed, which is what a rotation does — and a goroutine
// doing work in the background will call this in that window. Taking the
// whole process down for it, in a library whose job is to make a phone
// bearable, is not the right answer; an empty path that fails the next
// operation with a reason is.
func (a *App) Info() Info {
	if a == nil {
		return Info{}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.info
}

// drain keeps answering the UI thread after the app's function has returned.
// Without it, a program that simply stops would leave Android blocked in a
// callback and the app would be killed for not responding rather than
// closing.
func (a *App) drain() {
	for a.Alive() {
		if _, ok := a.Next(); !ok {
			return
		}
	}
	a.release()
}

var (
	mainMu   sync.Mutex
	mainFn   func(*App)
	fallback func(*App)
)

// Main registers the function that is the app. Call it from an init
// function: a shared library's main is never run, so init is the only place
// that happens before Android calls in.
//
//	func init() { app.Main(run) }
//
// Calling it twice replaces the first, which is only ever a mistake, so the
// second call says so in the log.
func Main(f func(*App)) {
	mainMu.Lock()
	defer mainMu.Unlock()
	if mainFn != nil {
		ndk.Warnf("app.Main called twice; the second function wins")
	}
	mainFn = f
}

// Run registers a function that is the whole app and takes no argument. It
// is for a program that reaches Android through antui.Window rather than
// through this package.
func Run(f func()) { Main(func(*App) { f() }) }

// Fallback registers f as the app only if nothing registers one with [Main].
//
// It is what the antuiapk command injects, pointing at the program's own
// main, so that a program written for a desktop runs on a phone with no
// Android in its source at all. Which init runs first does not matter: the
// choice is made when Android calls in, not when the registration happens.
func Fallback(f func()) {
	mainMu.Lock()
	defer mainMu.Unlock()
	fallback = func(*App) { f() }
}

// Current is the running activity, or nil when this process is not one —
// which is every desktop build, and an Android process before
// ANativeActivity_onCreate has been called.
func Current() *App { return current() }

func registeredMain() func(*App) {
	mainMu.Lock()
	defer mainMu.Unlock()
	if mainFn != nil {
		return mainFn
	}
	return fallback
}
