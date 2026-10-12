//go:build android

package location

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Watch reports where the device is, as it moves.
//
// f runs on the UI thread, because that is the thread the listener was
// registered from and the platform calls back on the caller's looper. Do not
// take long in it.
//
// **This costs battery**, and the two providers cost very differently: the
// satellites can take a few per cent an hour and the towers almost nothing.
// Stop it the moment the app stops caring.
func Watch(opt Options, f func(Fix)) (stop func(), err error) {
	ok, err := held()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNoPermission
	}
	providers := []string{opt.Provider}
	if opt.Provider == "" {
		providers = []string{GPS, Network}
	}

	token, release := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		if method != "onLocationChanged" || len(args) == 0 || args[0].IsNil() {
			return
		}
		if fix, err := read(e, args[0]); err == nil {
			f(fix)
		}
	})

	var listener jni.Object
	fail := func(err error) (func(), error) {
		release()
		return nil, err
	}
	if e := app.RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			manager, err := app.Service(e, app.ServiceLocation)
			if err != nil {
				return err
			}
			// LocationListener is an interface, so a proxy is enough and no
			// Java had to be written for it.
			l, err := app.Proxy(e, token, "android.location.LocationListener")
			if err != nil {
				return err
			}
			listener = e.Global(l)
			for _, p := range providers {
				jp, err := e.String(p)
				if err != nil {
					return err
				}
				if err := e.InvokeVoid(manager, "requestLocationUpdates",
					jni.Sig(jni.TVoid, jni.TString, jni.TLong, jni.TFloat,
						jni.TClass("android/location/LocationListener")),
					jni.Ref(jp), jni.Long(opt.Interval.Milliseconds()),
					jni.Float(opt.Distance), jni.Ref(listener)); err != nil {
					// A provider the device does not have throws rather than
					// reporting; with both asked for, one is enough.
					if len(providers) == 1 {
						return err
					}
				}
			}
			return nil
		})
	}); e != nil {
		return fail(e)
	}
	if err != nil {
		return fail(err)
	}

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
	}, nil
}
