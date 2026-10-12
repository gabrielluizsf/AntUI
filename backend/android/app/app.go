//go:build android

package app

import (
	"sync"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// App is the running activity. There is only ever one, and [Main] is given
// it when the system creates it.
type App struct {
	msgs chan message
	// dead is closed when the app goroutine will never receive again, so a
	// callback on the UI thread gives up instead of blocking forever.
	dead chan struct{}

	// Everything below is written on the UI thread and read on the app's,
	// so it is behind the mutex even when it looks like it could not be.
	mu     sync.Mutex
	window *ndk.Window
	rect   ndk.Rect
	focus  bool
	info   Info
	config Config
	state  State
	// saved is what a previous run of this activity wrote down; see
	// state.go.
	saved []byte

	// Everything below belongs to the app's goroutine alone and is never
	// touched from a callback, so none of it is behind the mutex.
	pending message
	held    bool
	done    bool
	// loop is this thread's looper, and queue the input attached to it.
	// Both are thread-local in the platform's sense, which is why the
	// goroutine is pinned to an OS thread.
	loop  *ndk.Looper
	queue *ndk.InputQueue

	// assets is the manager the platform handed the activity. It is set
	// once, before the app's goroutine starts, and never changes.
	assets unsafe.Pointer
}
