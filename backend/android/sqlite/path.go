//go:build android

package sqlite

import "github.com/gabrielluizsf/antui/backend/android/jni"

// Path is where the file is, for a caller that wants to copy or back it up.
func (db *DB) Path() (string, error) {
	var out string
	err := jni.Do(func(e *jni.Env) error {
		var err error
		out, err = e.InvokeString(db.obj, "getPath", jni.Sig(jni.TString))
		return err
	})
	return out, err
}
