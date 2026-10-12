//go:build android

package sqlite

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Close shuts the database. Rows opened from it must be closed first.
func (db *DB) Close() error {
	if db.closed {
		return nil
	}
	db.closed = true
	return jni.Do(func(e *jni.Env) error {
		e.InvokeVoid(db.obj, "close", jni.Sig(jni.TVoid))
		e.DeleteGlobal(db.obj)
		return nil
	})
}
