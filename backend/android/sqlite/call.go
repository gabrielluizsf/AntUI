//go:build android

package sqlite

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func (db *DB) call(method string) error {
	return jni.Do(func(e *jni.Env) error {
		return e.InvokeVoid(db.obj, method, jni.Sig(jni.TVoid))
	})
}
