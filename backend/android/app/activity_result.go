//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// ActivityResult is what another activity sent back.
type ActivityResult struct {
	// Code is Activity.RESULT_OK (-1) when the user finished, and
	// RESULT_CANCELED (0) when they backed out.
	Code int
	// Data is the Intent that came back, as a **global reference the
	// receiver owns**: it has to be given back with jni.Env.DeleteGlobal, or
	// the reference table fills up. It is nil when nothing came back, which
	// is the usual case for a cancel.
	Data jni.Object
}
