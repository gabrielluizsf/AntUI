//go:build android

package permission

import "github.com/gabrielluizsf/antui/backend/android/app"

// Request asks the user and waits.
//
// **Not from the frame loop.** The dialog pauses the app, Android waits for
// the app to acknowledge the pause, and a loop that is waiting here is a
// loop that is not acknowledging. Run it on a goroutine of its own.
func Request(names ...string) (Result, error) { return app.RequestPermissions(names...) }
