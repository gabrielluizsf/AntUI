//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Context is the app's android.content.Context — which is the activity, and
// is the starting point for almost everything Android offers. It is a global
// reference the library owns and callers must not free.
func Context() jni.Object { return jni.Activity() }
