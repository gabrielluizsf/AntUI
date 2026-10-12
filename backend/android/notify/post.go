//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Post puts a notification in the shade, or replaces the one with the same
// id.
//
// From API 33 this needs the POST_NOTIFICATIONS permission, and without it
// nothing appears and nothing fails. Ask for it with
// [antui/backend/android/permission.Ensure] before the first post.
func Post(n Notification) error {
	a := app.Current()
	if a == nil {
		return app.ErrNoUIThread
	}
	sdk := a.Config().SDK
	return jni.Do(func(e *jni.Env) error {
		builder, err := newBuilder(e, sdk, n.Channel)
		if err != nil {
			return err
		}
		if err := configureBuilder(e, builder, n, sdk); err != nil {
			return err
		}
		return builderPost(e, builder, n)
	})
}
