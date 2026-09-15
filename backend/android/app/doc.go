// Package app is the shape of an Android app: the activity Android created,
// the surface it draws on, and the lifecycle it is dragged through.
//
// Android does not run a native app the way a desktop does. There is no main
// to call: the system loads this library into a process it already made,
// calls one function in it, and from then on speaks in callbacks on a thread
// the app does not own. What this package does is turn that into something a
// Go program can read from top to bottom — a goroutine of the app's own, and
// a stream of events to take one at a time.
//
// # Registering
//
//	func init() { app.Main(run) }
//
//	func run(a *app.App) {
//		for {
//			e, ok := a.Next()
//			if !ok {
//				return
//			}
//			...
//		}
//	}
//
// It is init and not main because a shared library has no main: Go runs
// every package's init when the runtime starts and then stops, and the
// framework's call into [ANativeActivity_onCreate] is what starts it. The
// antuiapk command papers over this, so an ordinary main works too.
//
// # Which thread
//
// Every callback Android makes arrives on the UI thread, and blocking it is
// how an app gets killed for not responding. So callbacks do almost nothing:
// they hand the event to the app's own goroutine and, for the few that must,
// wait to be told it was dealt with. [App.Next] is where that waiting ends —
// see the note there, because getting it wrong is a crash rather than a bug.
package app
