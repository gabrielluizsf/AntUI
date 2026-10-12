//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// RequestPermissions asks the user and waits for the answer.
//
// **Do not call it from the frame loop.** Asking puts a dialog in front of
// the app, which pauses it — and Android waits for the app to acknowledge
// that pause. A frame loop that is waiting here is a frame loop that is not
// acknowledging, and both sides stop. Run it on a goroutine of its own and
// let the frame loop carry on drawing.
//
// A permission that is already granted comes back straight away, with no
// dialog. One the user has refused twice comes back refused, again with no
// dialog, and that is the platform's decision rather than this library's.
func RequestPermissions(names ...string) (PermissionResult, error) {
	if len(names) == 0 {
		return PermissionResult{}, nil
	}
	a := Current()
	if a == nil {
		return PermissionResult{}, ErrNoUIThread
	}
	if a.Config().SDK < 23 {
		granted := make([]bool, len(names))
		for i := range granted {
			granted[i] = true
		}
		return PermissionResult{Permissions: names, Granted: granted}, nil
	}

	code := nextRequestCode()
	ch := awaitPermission(code)
	var callErr error
	if err := RunOnUISync(func() {
		callErr = jni.Do(func(e *jni.Env) error {
			act := jni.Activity()
			cls, err := e.ClassOf(act)
			if err != nil {
				return err
			}
			defer e.Delete(cls.Object())
			m, err := e.Method(cls, "requestPermissions",
				jni.Sig(jni.TVoid, jni.TArray(jni.TString), jni.TInt))
			if err != nil {
				return err
			}
			arr, err := e.NewStrings(names)
			if err != nil {
				return err
			}
			return e.CallVoid(act, m, jni.Ref(arr), jni.Int(code))
		})
	}); err != nil {
		callErr = err
	}
	if callErr != nil {
		waitMu.Lock()
		delete(permWait, code)
		waitMu.Unlock()
		return PermissionResult{}, callErr
	}
	return <-ch, nil
}
