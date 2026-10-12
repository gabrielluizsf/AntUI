//go:build android

package sqlite

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Next moves to the next row and reports whether there was one.
func (r *Rows) Next() bool {
	if r.closed || r.err != nil {
		return false
	}
	var more bool
	if err := jni.Do(func(e *jni.Env) error {
		var err error
		more, err = e.InvokeBool(r.cursor, "moveToNext", jni.Sig(jni.TBool))
		return err
	}); err != nil {
		r.err = err
		return false
	}
	return more
}
