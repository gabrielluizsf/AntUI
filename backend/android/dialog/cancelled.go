//go:build android

package dialog

const (
	// Cancelled is the back button or a tap outside, which is not one of the
	// buttons and is the answer most code forgets to handle.
	Cancelled Button = 0
	// Positive is the one on the right, which is the one that does the
	// thing: OK, Save, Delete.
	Positive Button = -1
	// Negative is the one that does not: Cancel.
	Negative Button = -2
	// Neutral is the third, for an answer that is neither — "Not now",
	// "Tell me more".
	Neutral Button = -3
)
