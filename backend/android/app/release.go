//go:build android

package app

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
