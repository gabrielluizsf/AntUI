//go:build android

package notify

// OnAction says what a named action does.
//
// Register every one of them at startup, before anything is posted, and
// register them again on every launch — the tap that brings the app back may
// be the first thing that happens in a new process.
//
// f runs on a thread of the system's choosing, not the app's. reply is what
// the user typed, and is empty for a button that is only a button.
func OnAction(name string, f func(id int32, reply string)) {
	actionMu.Lock()
	actionFns[name] = f
	actionMu.Unlock()
	actionOnce.Do(startReceiver)
}
