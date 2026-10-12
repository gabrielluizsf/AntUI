//go:build android

package network

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// watchHandler reports the whole state on any change.
//
// Every one of the three callbacks means the same thing to a caller:
// something changed, here is how it is now. Reading it again costs a few
// calls and avoids piecing the state together from three different callbacks
// that each carry part of it.
func watchHandler(f func(State)) app.Handler {
	return func(e *jni.Env, method string, args []jni.Object) {
		if s, err := Read(); err == nil {
			f(s)
		}
	}
}
