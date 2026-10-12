package apk

// A Problem is one thing the store will act on.
type Problem struct {
	// Fatal is whether the store refuses the upload outright. A problem that
	// is not fatal is one it accepts and complains about, or one that will
	// cost the app users without anyone saying so.
	Fatal bool
	// Text is what is wrong, and Fix is what to do about it.
	Text string
	Fix  string
}
