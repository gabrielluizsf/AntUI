//go:build android

package network

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// watchStop returns the function that unregisters the callback, safe to call
// more than once.
func watchStop(callback jni.Object, release func()) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			jni.Do(func(e *jni.Env) error {
				manager, err := app.Service(e, app.ServiceConnectivity)
				if err == nil {
					e.InvokeVoid(manager, "unregisterNetworkCallback",
						jni.Sig(jni.TVoid,
							jni.TClass("android/net/ConnectivityManager$NetworkCallback")),
						jni.Ref(callback))
				}
				e.DeleteGlobal(callback)
				return nil
			})
			release()
		})
	}
}
