//go:build android

package location

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

// The two providers a device without Play Services has. There is no third:
// the "fused" provider everyone knows is part of Play Services, not of
// Android, and an app that must work without them uses these.
const (
	// GPS is the satellites: accurate to a few metres, slow to get a first
	// answer, useless indoors, and expensive in battery.
	GPS = "gps"
	// Network is the towers and the wifi around: accurate to tens or
	// hundreds of metres, answers immediately, works indoors, costs almost
	// nothing.
	Network = "network"
)

// ErrNoPermission is asking for a position without having been granted one
// of the two location permissions.
var ErrNoPermission = errors.New("location: this app has neither " +
	"ACCESS_FINE_LOCATION nor ACCESS_COARSE_LOCATION; ask for one with " +
	"antui/backend/android/permission before asking where the device is")

// Fix is where the device was, and how well it knew.
type Fix struct {
	Latitude, Longitude float64
	// Accuracy is the radius in metres inside which the device believes it
	// is, with about 68 per cent confidence. It is not a promise and it is
	// not a maximum.
	Accuracy float32
	// Altitude is metres above the WGS84 ellipsoid, which is not sea level
	// and can differ from it by a hundred metres. Zero when unknown.
	Altitude float64
	// Speed is metres a second and Bearing is degrees clockwise from north.
	// Both are zero on a device standing still and on one that cannot say.
	Speed   float32
	Bearing float32
	// At is when the position was taken.
	At time.Time
	// Provider is which of the two above answered.
	Provider string
}

// String is the position and how good it is, in one line.
func (f Fix) String() string {
	return fmt.Sprintf("%.6f, %.6f ±%.0fm (%s)", f.Latitude, f.Longitude, f.Accuracy, f.Provider)
}

// held reports whether either location permission has been granted.
func held() (bool, error) {
	for _, p := range []string{permission.FineLocation, permission.CoarseLocation} {
		ok, err := permission.Held(p)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// Last is the most recent position the system already has, without turning
// anything on.
//
// It is what almost every app should ask first: it is instant and free, and
// on a phone that has been outdoors in the last few minutes it is as good as
// anything a new fix would give. It reports false when the system has
// nothing — a phone just switched on, or one that has never been given the
// permission.
func Last() (Fix, bool, error) {
	ok, err := held()
	if err != nil {
		return Fix{}, false, err
	}
	if !ok {
		return Fix{}, false, ErrNoPermission
	}

	var best Fix
	found := false
	err = jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceLocation)
		if err != nil {
			return err
		}
		for _, provider := range []string{GPS, Network} {
			jp, err := e.String(provider)
			if err != nil {
				return err
			}
			loc, err := e.Invoke(manager, "getLastKnownLocation",
				jni.Sig(jni.TClass("android/location/Location"), jni.TString), jni.Ref(jp))
			if err != nil || loc.IsNil() {
				continue
			}
			fix, err := read(e, loc)
			if err != nil {
				continue
			}
			// The newer of the two, not the more accurate: a stale reading
			// from the satellites is worse than a fresh one from the towers,
			// however tight its stated accuracy.
			if !found || fix.At.After(best.At) {
				best, found = fix, true
			}
		}
		return nil
	})
	return best, found, err
}

// Options is how a watch is opened.
type Options struct {
	// Provider is [GPS] or [Network]. Empty asks for both, and reports from
	// whichever answers.
	Provider string
	// Interval is the shortest gap between reports. It is a floor the
	// platform is allowed to exceed and, on API 29 and up, to ignore
	// downwards — the system throttles background apps whatever they ask.
	Interval time.Duration
	// Distance is how far the device must move before it is worth another
	// report, in metres. Zero reports on time alone.
	Distance float32
}

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

// Enabled reports whether a provider is switched on. The permission being
// granted is not enough: the user can turn location off for the whole
// device, and then every request answers nothing for ever.
func Enabled(provider string) (bool, error) {
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceLocation)
		if err != nil {
			return err
		}
		jp, err := e.String(provider)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(manager, "isProviderEnabled",
			jni.Sig(jni.TBool, jni.TString), jni.Ref(jp))
		return err
	})
	return out, err
}

// read pulls a Fix out of an android.location.Location.
func read(e *jni.Env, loc jni.Object) (Fix, error) {
	var f Fix
	var err error
	if f.Latitude, err = invokeDouble(e, loc, "getLatitude"); err != nil {
		return f, err
	}
	if f.Longitude, err = invokeDouble(e, loc, "getLongitude"); err != nil {
		return f, err
	}
	if f.Altitude, err = invokeDouble(e, loc, "getAltitude"); err != nil {
		return f, err
	}
	if f.Accuracy, err = e.InvokeFloat(loc, "getAccuracy", jni.Sig(jni.TFloat)); err != nil {
		return f, err
	}
	if f.Speed, err = e.InvokeFloat(loc, "getSpeed", jni.Sig(jni.TFloat)); err != nil {
		return f, err
	}
	if f.Bearing, err = e.InvokeFloat(loc, "getBearing", jni.Sig(jni.TFloat)); err != nil {
		return f, err
	}
	ms, err := e.InvokeLong(loc, "getTime", jni.Sig(jni.TLong))
	if err != nil {
		return f, err
	}
	f.At = time.UnixMilli(ms)
	if f.Provider, err = e.InvokeString(loc, "getProvider", jni.Sig(jni.TString)); err != nil {
		return f, err
	}
	return f, nil
}
