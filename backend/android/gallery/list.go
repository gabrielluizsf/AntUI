//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// List is what is in one of the device's libraries.
//
// **Seeing other apps' media needs a permission** — READ_EXTERNAL_STORAGE
// below Android 13, and READ_MEDIA_IMAGES, _VIDEO or _AUDIO from 13 — and
// without one the answer is only what this app put there, with no error to
// say so. [Query.Mine] asks for that on purpose, and is what an app that
// only wants its own should use: it needs nothing granted.
func List(kind Kind, q Query) ([]Item, error) {
	a := app.Current()
	if a == nil {
		return nil, app.ErrNoUIThread
	}
	modern := a.Config().SDK >= 29

	var out []Item
	err := jni.Do(func(e *jni.Env) error {
		columns := listColumns
		if modern {
			columns = append(append([]string{}, columns...), colPath)
		}
		projection, err := e.NewStrings(columns)
		if err != nil {
			return err
		}

		var where jni.Object
		var args jni.Object
		switch {
		case q.Mine:
			// Everything MediaStore reports to an app with no permission is
			// already its own, but saying so makes the intent plain and
			// keeps the answer the same when a permission is later added.
			if where, err = e.String("owner_package_name = ?"); err != nil {
				return err
			}
			name, err := e.InvokeString(app.Context(), "getPackageName", jni.Sig(jni.TString))
			if err != nil {
				return err
			}
			if args, err = e.NewStrings([]string{name}); err != nil {
				return err
			}
		case q.Album != "" && modern:
			if where, err = e.String(colPath + " LIKE ?"); err != nil {
				return err
			}
			if args, err = e.NewStrings([]string{"%/" + q.Album + "/%"}); err != nil {
				return err
			}
		}
		// Newest first, which is what a gallery shows and what a limit
		// should keep.
		order, err := e.String("date_added DESC")
		if err != nil {
			return err
		}

		class, field := kind.collection()
		collection, err := e.Constant(class, field, jni.TClass("android/net/Uri"))
		if err != nil {
			return err
		}
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		cursor, err := e.Invoke(resolver, "query",
			jni.Sig(jni.TClass("android/database/Cursor"),
				jni.TClass("android/net/Uri"), jni.TArray(jni.TString),
				jni.TString, jni.TArray(jni.TString), jni.TString),
			jni.Ref(collection), jni.Ref(projection), jni.Ref(where),
			jni.Ref(args), jni.Ref(order))
		if err != nil || cursor.IsNil() {
			return err
		}
		defer e.InvokeVoid(cursor, "close", jni.Sig(jni.TVoid))

		for {
			more, err := e.InvokeBool(cursor, "moveToNext", jni.Sig(jni.TBool))
			if err != nil || !more {
				return err
			}
			item, err := readItem(e, cursor, collection, modern)
			if err != nil {
				return err
			}
			out = append(out, item)
			if q.Limit > 0 && len(out) >= q.Limit {
				return nil
			}
		}
	})
	return out, err
}
