//go:build android

package sqlite

import "github.com/gabrielluizsf/antui/backend/android/jni"

// DB is an open database.
//
// SQLite is on every Android device and has been since the first one, but
// only through Java: the C library is there and the NDK does not expose it,
// and linking a copy of SQLite into the app would mean shipping a second one
// beside the platform's and having two caches of the same file.
type DB struct {
	obj    jni.Object
	closed bool
}
