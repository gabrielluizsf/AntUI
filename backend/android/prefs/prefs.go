//go:build android

package prefs

// Prefs is one named store of small values, kept as XML in the app's own
// directory and read back on the next run.
//
// It is for settings and for the scraps of state a program wants to
// remember: which level was reached, whether the sound is on, the last
// window it had. It is not a database and not a file — a value larger than a
// few kilobytes belongs somewhere else, and the whole store is loaded into
// memory the first time it is opened.
type Prefs struct {
	name string
}
