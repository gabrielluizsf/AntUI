//go:build android

package app

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
