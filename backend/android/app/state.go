//go:build android

package app

// State is where the activity is in its life. Android's own diagram is a
// cycle rather than a line: an app goes Resumed → Paused → Stopped and back
// up again every time the user switches away and returns, and only reaches
// Destroyed once.
type State int
