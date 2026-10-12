//go:build android

package network

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Watch calls f whenever the connection changes — a network appearing, the
// one in use going away, or what it can do changing, which is how a
// connection going from metered to unmetered is noticed.
//
// f runs on a thread of the system's choosing, not the app's, and is given
// the whole state rather than the difference, because the difference is
// almost never what a caller wants. Call stop when done; a callback left
// registered outlives the app's interest in it and goes on being called.
func Watch(f func(State)) (stop func(), err error) {
	token, release := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		// Every one of the three means the same thing to a caller: something
		// changed, here is how it is now. Reading it again costs a few calls
		// and avoids having to piece the state together from three
		// different callbacks that each carry part of it.
		if s, err := Read(); err == nil {
			f(s)
		}
	})

	var callback jni.Object
	err = jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceConnectivity)
		if err != nil {
			return err
		}
		cb, err := app.Listener(e, app.NetworkCallbackClass, token)
		if err != nil {
			return err
		}
		// A global reference, and not only so it can be unregistered: the
		// only thing holding this object is the platform, and a local
		// reference would be gone before the first callback.
		callback = e.Global(cb)

		builder, err := e.Make("android/net/NetworkRequest$Builder", jni.Sig(jni.TVoid))
		if err != nil {
			return err
		}
		if _, err := e.Invoke(builder, "addCapability",
			jni.Sig(jni.TClass("android/net/NetworkRequest$Builder"), jni.TInt),
			jni.Int(capInternet)); err != nil {
			return err
		}
		request, err := e.Invoke(builder, "build",
			jni.Sig(jni.TClass("android/net/NetworkRequest")))
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "registerNetworkCallback",
			jni.Sig(jni.TVoid, jni.TClass("android/net/NetworkRequest"),
				jni.TClass("android/net/ConnectivityManager$NetworkCallback")),
			jni.Ref(request), jni.Ref(callback))
	})
	if err != nil {
		release()
		return func() {}, err
	}

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
	}, nil
}
