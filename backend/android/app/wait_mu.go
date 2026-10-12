//go:build android

package app

import (
	"sync"
	"sync/atomic"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// The registry. A request goes out with a code, and the callback that comes
// back carries the same code — which is the only way to tell one answer from
// another when two are outstanding.
var (
	waitMu   sync.Mutex
	permWait = map[int32]chan PermissionResult{}
	actWait  = map[int32]chan ActivityResult{}
	lastCode atomic.Int32

	intentMu sync.Mutex
	intentFn func(jni.Object)
)
