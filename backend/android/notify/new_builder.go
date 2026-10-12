//go:build android

package notify

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// newBuilder makes the right Notification.Builder for this platform: the one
// that takes a channel from API 26, and the one that does not below it.
func newBuilder(e *jni.Env, sdk int, channel string) (jni.Object, error) {
	if sdk >= 26 {
		if channel == "" {
			return jni.Object{}, fmt.Errorf(
				"notify: this device needs a channel, and none was named")
		}
		id, err := e.String(channel)
		if err != nil {
			return jni.Object{}, err
		}
		return e.Make("android/app/Notification$Builder",
			jni.Sig(jni.TVoid, jni.TContext, jni.TString),
			jni.Ref(app.Context()), jni.Ref(id))
	}
	return e.Make("android/app/Notification$Builder",
		jni.Sig(jni.TVoid, jni.TContext), jni.Ref(app.Context()))
}
