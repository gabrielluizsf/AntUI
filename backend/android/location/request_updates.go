//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// requestUpdates registers the listener with every provider asked for and
// returns the global reference to it.
func requestUpdates(e *jni.Env, token int64, opt Options, providers []string) (jni.Object, error) {
	manager, err := app.Service(e, app.ServiceLocation)
	if err != nil {
		return jni.Object{}, err
	}
	// LocationListener is an interface, so a proxy is enough and no Java had
	// to be written for it.
	l, err := app.Proxy(e, token, "android.location.LocationListener")
	if err != nil {
		return jni.Object{}, err
	}
	listener := e.Global(l)
	for _, p := range providers {
		jp, err := e.String(p)
		if err != nil {
			return jni.Object{}, err
		}
		if err := e.InvokeVoid(manager, "requestLocationUpdates",
			jni.Sig(jni.TVoid, jni.TString, jni.TLong, jni.TFloat,
				jni.TClass("android/location/LocationListener")),
			jni.Ref(jp), jni.Long(opt.Interval.Milliseconds()),
			jni.Float(opt.Distance), jni.Ref(listener)); err != nil {
			// A provider the device does not have throws rather than
			// reporting; with both asked for, one is enough.
			if len(providers) == 1 {
				return jni.Object{}, err
			}
		}
	}
	return listener, nil
}
