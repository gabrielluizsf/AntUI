//go:build android

package audio

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// unpluggedStop returns the function that unregisters the receiver, safe to
// call more than once.
func unpluggedStop(receiver jni.Object, releaseToken func()) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			jni.Do(func(e *jni.Env) error {
				e.InvokeVoid(app.Context(), "unregisterReceiver",
					jni.Sig(jni.TVoid, jni.TClass("android/content/BroadcastReceiver")),
					jni.Ref(receiver))
				e.DeleteGlobal(receiver)
				return nil
			})
			releaseToken()
		})
	}
}
