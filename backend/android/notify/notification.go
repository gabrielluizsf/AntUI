//go:build android

package notify

// Notification is one message in the shade.
type Notification struct {
	// ID identifies it within this app. Posting twice with the same one
	// replaces rather than adds, which is how a progress bar moves.
	ID int32
	// Channel is the id of the channel it belongs to.
	Channel string
	Title   string
	Text    string

	// Icon is a drawable resource id. Zero means the platform's own
	// information icon, which is there on every device — an app has no
	// resources of its own until it is built with some.
	Icon int32

	// Progress draws a bar. Total is the top of the range; a Total of 0 means
	// no bar, and Indeterminate makes it a moving stripe with no position.
	Progress      int32
	Total         int32
	Indeterminate bool

	// Ongoing marks it as something in progress: it cannot be swiped away.
	Ongoing bool
	// AutoCancel dismisses it when the user taps it, which is what almost
	// every notification should do.
	AutoCancel bool
	// Silent posts without a sound even on a channel that makes one.
	Silent bool

	// Actions are the buttons under the notification. See [Action]: each
	// carries a name rather than a function, and [OnAction] is where the
	// names are given meaning.
	Actions []Action

	// Extra is handed back through [antui/backend/android/app.OnNewIntent] when the
	// user taps it, which is how an app knows *which* notification was
	// tapped.
	Extra map[string]string
}
