//go:build android

package antui

import "errors"

// errNotAndroidApp is what an Android *process* gets when it is not an
// Android *app* — a command-line binary pushed to a device with adb, say,
// which has a Go runtime and no activity.
var errNotAndroidApp = errors.New("antui: no Android activity: this process was not started by the system as an app")
