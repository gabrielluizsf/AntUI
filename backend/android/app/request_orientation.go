//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/jni"

func requestOrientation(want int) error {
	var err error
	if e := RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			act := jni.Activity()
			cls, err := e.ClassOf(act)
			if err != nil {
				return err
			}
			defer e.Delete(cls.Object())
			m, err := e.Method(cls, "setRequestedOrientation", jni.Sig(jni.TVoid, jni.TInt))
			if err != nil {
				return err
			}
			return e.CallVoid(act, m, jni.Int(int32(want)))
		})
	}); e != nil {
		return e
	}
	return err
}
