//go:build android

package sqlite

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// boxAll turns Go values into the Objects execSQL binds.
func boxAll(e *jni.Env, args []any) (jni.Object, error) {
	cls, err := e.Class("java/lang/Object")
	if err != nil {
		return jni.Object{}, err
	}
	arr, err := e.NewObjects(cls, len(args))
	if err != nil {
		return jni.Object{}, err
	}
	err = e.Frame(len(args)*4+16, func() error {
		for i, a := range args {
			boxed, err := box(e, a)
			if err != nil {
				return err
			}
			if err := e.SetIndex(arr, i, boxed); err != nil {
				return err
			}
		}
		return nil
	})
	return arr, err
}
