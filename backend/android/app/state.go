//go:build android

package app

/*
#include <stdlib.h>
#include "glue.h"
*/
import "C"

import (
	"sync"
	"unsafe"
)

// State is where the activity is in its life. Android's own diagram is a
// cycle rather than a line: an app goes Resumed → Paused → Stopped and back
// up again every time the user switches away and returns, and only reaches
// Destroyed once.
type State int

// The states, in the order they are first reached.
const (
	// Created is after onCreate and before onStart. There is no window yet.
	Created State = iota
	// Started means the activity is visible but not in front — behind a
	// dialog, or beside another app on a split screen.
	Started
	// Resumed means the user is looking at it and it has the input. This is
	// the only state in which an app should be doing work.
	Resumed
	// Paused is the moment after losing the front. **Anything that must
	// survive has to be saved by the time this returns**: the process may be
	// frozen or killed with no further warning.
	Paused
	// Stopped means no longer visible at all.
	Stopped
	// Destroyed is the end.
	Destroyed
)

var stateNames = [...]string{"created", "started", "resumed", "paused", "stopped", "destroyed"}

// String names the state, which is what a log line wants.
func (s State) String() string {
	if int(s) < len(stateNames) {
		return stateNames[s]
	}
	return "state(?)"
}

// State is where the activity is now.
func (a *App) State() State {
	if a == nil {
		return Destroyed
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state
}

// Foreground reports whether the user is looking at the app. It is the one
// question most programs actually want to ask.
func (a *App) Foreground() bool { return a.State() == Resumed }

// noteState folds a lifecycle event into the state. It is called from apply,
// on the app's goroutine, before the event is handed to the caller — so an
// app that asks the state while handling Pause is told Paused.
func (a *App) noteState(e Event) {
	switch e {
	case Start:
		a.state = Started
	case Resume:
		a.state = Resumed
	case Pause:
		a.state = Paused
	case Stop:
		a.state = Stopped
	case Destroy:
		a.state = Destroyed
	}
}

// Saved state.

var (
	saveMu sync.Mutex
	saveFn func() []byte
)

// SaveState registers what to write down when Android asks. It is asked
// around the time the app is paused, and what comes back is handed to the
// next launch through [App.RestoredState] — but only if the process was
// killed while in the background. A user who closes the app deliberately
// gets nothing back, which is the point.
//
// **f is called on the UI thread**, synchronously, and may not wait for the
// app's own goroutine — the UI thread is what would have to run first for
// that to finish. Whatever it reads has to be safe to read from there.
//
// Keep it small. It travels through the system as part of a transaction with
// a hard size limit, and going over it kills the app rather than truncating.
func SaveState(f func() []byte) {
	saveMu.Lock()
	saveFn = f
	saveMu.Unlock()
}

// RestoredState is what a previous run of this activity saved, or nil.
func (a *App) RestoredState() []byte {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.saved
}

// antuiOnSaveInstanceState is called by the framework, on the UI thread,
// when it wants the app's state written down.
//
// The block handed back is freed by the framework with free, so it has to
// come from malloc — a Go slice's memory would be freed by something that
// does not own it, which is a corruption rather than a crash.
//
//export antuiOnSaveInstanceState
func antuiOnSaveInstanceState(_ *C.ANativeActivity, size *C.size_t) unsafe.Pointer {
	*size = 0
	saveMu.Lock()
	f := saveFn
	saveMu.Unlock()
	if f == nil {
		return nil
	}
	b := f()
	if len(b) == 0 {
		return nil
	}
	p := C.malloc(C.size_t(len(b)))
	if p == nil {
		return nil
	}
	copy(unsafe.Slice((*byte)(p), len(b)), b)
	*size = C.size_t(len(b))
	return p
}
