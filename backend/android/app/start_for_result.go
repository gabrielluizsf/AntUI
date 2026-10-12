//go:build android

package app

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// StartForResult starts another activity and waits for what it sends back —
// a photo from the camera, a file from the picker, a document from a cloud
// drive.
//
// The same warning as [RequestPermissions] applies, and more strongly: the
// other activity is in front for as long as the user wants it to be.
//
// The Intent in the result is a global reference the caller owns.
func StartForResult(intent jni.Object) (ActivityResult, error) {
	if intent.IsNil() {
		return ActivityResult{}, errors.New("app: there is no intent to start")
	}
	code := nextRequestCode()
	ch := awaitActivity(code)
	var callErr error
	if err := RunOnUISync(func() {
		callErr = jni.Do(func(e *jni.Env) error {
			act := jni.Activity()
			cls, err := e.ClassOf(act)
			if err != nil {
				return err
			}
			defer e.Delete(cls.Object())
			m, err := e.Method(cls, "startActivityForResult",
				jni.Sig(jni.TVoid, jni.TClass("android/content/Intent"), jni.TInt))
			if err != nil {
				return err
			}
			return e.CallVoid(act, m, jni.Ref(intent), jni.Int(code))
		})
	}); err != nil {
		callErr = err
	}
	if callErr != nil {
		waitMu.Lock()
		delete(actWait, code)
		waitMu.Unlock()
		return ActivityResult{}, callErr
	}
	return <-ch, nil
}
