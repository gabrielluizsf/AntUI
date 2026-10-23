//go:build android

package antui

import (
	"sync/atomic"

	"github.com/gabrielluizsf/antui/backend/android/notify"
)

// notificationChannel is the one channel SendNotification posts to. A
// notification outside a created channel is silently dropped from API 26,
// so it is created before every post — which costs nothing once it exists,
// and survives the user deleting it in the settings.
var notificationChannel = notify.Channel{ID: "antui", Name: "Notifications", Importance: notify.Default}

// notificationID numbers the posts, so that each is a notification of its
// own rather than a replacement of the last.
var notificationID atomic.Int32

// postNotification creates the channel and posts title and body to it.
func postNotification(title, body string) error {
	if err := notify.CreateChannel(notificationChannel); err != nil {
		return err
	}
	return notify.Post(notify.Notification{
		ID:         notificationID.Add(1),
		Channel:    notificationChannel.ID,
		Title:      title,
		Text:       body,
		AutoCancel: true,
	})
}
