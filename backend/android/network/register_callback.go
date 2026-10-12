//go:build android

package network

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// registerCallback registers for the internet capability and returns the
// global reference to the callback.
func registerCallback(e *jni.Env, token int64) (jni.Object, error) {
	manager, err := app.Service(e, app.ServiceConnectivity)
	if err != nil {
		return jni.Object{}, err
	}
	cb, err := app.Listener(e, app.NetworkCallbackClass, token)
	if err != nil {
		return jni.Object{}, err
	}
	// A global reference, and not only so it can be unregistered: the only
	// thing holding this object is the platform, and a local reference would
	// be gone before the first callback.
	callback := e.Global(cb)

	builder, err := e.Make("android/net/NetworkRequest$Builder", jni.Sig(jni.TVoid))
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(builder, "addCapability",
		jni.Sig(jni.TClass("android/net/NetworkRequest$Builder"), jni.TInt),
		jni.Int(capInternet)); err != nil {
		return jni.Object{}, err
	}
	request, err := e.Invoke(builder, "build",
		jni.Sig(jni.TClass("android/net/NetworkRequest")))
	if err != nil {
		return jni.Object{}, err
	}
	if err := e.InvokeVoid(manager, "registerNetworkCallback",
		jni.Sig(jni.TVoid, jni.TClass("android/net/NetworkRequest"),
			jni.TClass("android/net/ConnectivityManager$NetworkCallback")),
		jni.Ref(request), jni.Ref(callback)); err != nil {
		return jni.Object{}, err
	}
	return callback, nil
}
