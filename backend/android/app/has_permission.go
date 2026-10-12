//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// HasPermission reports whether a permission has already been granted.
//
// Below API 23 every permission an app declares is granted when it is
// installed, and there is nothing to ask, so this is always true there.
func HasPermission(name string) (bool, error) {
	a := Current()
	if a == nil {
		return false, ErrNoUIThread
	}
	if a.Config().SDK < 23 {
		return true, nil
	}
	var granted bool
	err := jni.Do(func(e *jni.Env) error {
		act := jni.Activity()
		cls, err := e.ClassOf(act)
		if err != nil {
			return err
		}
		defer e.Delete(cls.Object())
		m, err := e.Method(cls, "checkSelfPermission", jni.Sig(jni.TInt, jni.TString))
		if err != nil {
			return err
		}
		js, err := e.String(name)
		if err != nil {
			return err
		}
		r, err := e.CallInt(act, m, jni.Ref(js))
		granted = r == 0
		return err
	})
	return granted, err
}
