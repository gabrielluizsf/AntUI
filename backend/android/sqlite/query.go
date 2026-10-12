//go:build android

package sqlite

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Query runs a SELECT.
//
// **Every argument is bound as text.** That is the platform's API and not a
// choice here: rawQuery takes a String array and nothing else. SQLite's own
// typing means "42" compares equal to 42 in almost every case, but a column
// declared with no affinity will hold the text — so a number that must stay
// a number belongs in the statement or in an Exec.
func (db *DB) Query(sql string, args ...any) (*Rows, error) {
	if db.closed {
		return nil, errors.New("sqlite: the database is closed")
	}
	rows := &Rows{}
	err := jni.Do(func(e *jni.Env) error {
		jsql, err := e.String(sql)
		if err != nil {
			return err
		}
		var jargs jni.Object
		if len(args) > 0 {
			text := make([]string, len(args))
			for i, a := range args {
				text[i] = asText(a)
			}
			if jargs, err = e.NewStrings(text); err != nil {
				return err
			}
		}
		cursor, err := e.Invoke(db.obj, "rawQuery",
			jni.Sig(jni.TClass("android/database/Cursor"), jni.TString,
				jni.TArray(jni.TString)),
			jni.Ref(jsql), jni.Ref(jargs))
		if err != nil {
			return err
		}
		if cursor.IsNil() {
			return errors.New("sqlite: the query answered nothing at all")
		}
		rows.cursor = e.Global(cursor)
		names, err := e.Invoke(rows.cursor, "getColumnNames",
			jni.Sig(jni.TArray(jni.TString)))
		if err == nil && !names.IsNil() {
			rows.columns, _ = e.GoStrings(names)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}
