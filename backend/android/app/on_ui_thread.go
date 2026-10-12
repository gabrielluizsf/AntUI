//go:build android

package app

/*
#include <android/looper.h>
#include "glue.h"
*/
import "C"

// OnUIThread reports whether the caller is already on the UI thread. It is
// what stops [RunOnUISync] deadlocking on itself.
func OnUIThread() bool {
	return uiLooper != nil && C.ALooper_forThread() == uiLooper
}
