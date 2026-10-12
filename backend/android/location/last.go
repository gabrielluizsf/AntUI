//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
