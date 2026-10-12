//go:build android

package location

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// watchStop returns the function that stops the updates, safe to call more
// than once.
func watchStop(listener jni.Object, release func()) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			app.RunOnUISync(func() {
				jni.Do(func(e *jni.Env) error {
					if manager, err := app.Service(e, app.ServiceLocation); err == nil {
						e.InvokeVoid(manager, "removeUpdates",
							jni.Sig(jni.TVoid, jni.TClass("android/location/LocationListener")),
							jni.Ref(listener))
					}
					e.DeleteGlobal(listener)
					return nil
				})
			})
			release()
		})
	}
}
