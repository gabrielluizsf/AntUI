//go:build android

package network

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Read asks about the current connection.
//
// It needs android.permission.ACCESS_NETWORK_STATE in the manifest, which is
// a normal permission — granted at install, with nothing to ask the user.
// Without it every call here reports a SecurityException.
func Read() (State, error) {
	out := State{Transport: None}
	a := app.Current()
	if a == nil {
		return out, app.ErrNoUIThread
	}
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceConnectivity)
		if err != nil {
			return err
		}
		if a.Config().SDK < 23 {
			return readOld(e, manager, &out)
		}
		return readModern(e, manager, &out)
	})
	return out, err
}
