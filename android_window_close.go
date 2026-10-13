//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/app"

// close asks Android to finish the activity. The process is not ended here:
// that is the system's to do, and a native app that calls exit is a native
// app that loses whatever the platform was about to save.
func (b *androidWindow) close() {
	b.surface = nil
	app.Finish()
}
