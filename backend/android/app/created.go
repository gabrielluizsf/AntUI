//go:build android

package app

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
