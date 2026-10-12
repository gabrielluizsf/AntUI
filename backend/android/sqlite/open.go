//go:build android

package sqlite

import (
	"errors"
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Open opens or creates a database in the app's own directory, by name.
// It is the ordinary case: the file lives where the system backs it up and
// deletes it with the app.
func Open(name string) (*DB, error) {
	if name == "" {
		return nil, errors.New("sqlite: a database needs a name")
	}
	var db DB
	err := jni.Do(func(e *jni.Env) error {
		jname, err := e.String(name)
		if err != nil {
			return err
		}
		obj, err := e.Invoke(app.Context(), "openOrCreateDatabase",
			jni.Sig(jni.TClass("android/database/sqlite/SQLiteDatabase"),
				jni.TString, jni.TInt,
				jni.TClass("android/database/sqlite/SQLiteDatabase$CursorFactory")),
			jni.Ref(jname), jni.Int(modePrivate), jni.Ref(jni.Object{}))
		if err != nil {
			return err
		}
		if obj.IsNil() {
			return fmt.Errorf("sqlite: cannot open %s", name)
		}
		db.obj = e.Global(obj)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &db, nil
}
