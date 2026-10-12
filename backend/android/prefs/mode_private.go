//go:build android

package prefs

// modePrivate is Context.MODE_PRIVATE: readable by this app and no other.
// It is the only mode still allowed — the world-readable ones were removed
// in Android 7 and throw if asked for.
const modePrivate = 0
