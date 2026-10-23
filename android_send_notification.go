//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

// sendNotification puts a notification in the shade.
//
// From API 33 posting needs POST_NOTIFICATIONS, which the app's manifest has
// to declare and the user has to grant; without it nothing appears and
// nothing fails. When it is not held yet, asking is a dialog that pauses the
// app — and waiting for that from the frame loop hangs it — so the request
// and the post that follows it go to a goroutine, and the answer here is
// only that the notification was asked for.
func (b *androidWindow) sendNotification(title, body string) bool {
	a := app.Current()
	if a == nil {
		return false
	}
	if a.Config().SDK >= 33 {
		held, err := permission.Held(permission.Notifications)
		if err != nil {
			return false
		}
		if !held {
			go func() {
				if granted, err := permission.Ensure(permission.Notifications); err == nil && granted {
					postNotification(title, body)
				}
			}()
			return true
		}
	}
	return postNotification(title, body) == nil
}
