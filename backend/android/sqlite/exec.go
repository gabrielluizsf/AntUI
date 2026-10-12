//go:build android

package sqlite

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Exec runs a statement that answers nothing: a CREATE, an INSERT, an
// UPDATE. The arguments replace the ? marks, in order.
func (db *DB) Exec(sql string, args ...any) error {
	if db.closed {
		return errors.New("sqlite: the database is closed")
	}
	return jni.Do(func(e *jni.Env) error {
		jsql, err := e.String(sql)
		if err != nil {
			return err
		}
		if len(args) == 0 {
			return e.InvokeVoid(db.obj, "execSQL",
				jni.Sig(jni.TVoid, jni.TString), jni.Ref(jsql))
		}
		boxed, err := boxAll(e, args)
		if err != nil {
			return err
		}
		return e.InvokeVoid(db.obj, "execSQL",
			jni.Sig(jni.TVoid, jni.TString, jni.TArray(jni.TObject)),
			jni.Ref(jsql), jni.Ref(boxed))
	})
}
