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
		where, args, err := listFilter(e, q, modern)
		if err != nil {
			return err
		}
		// Newest first, which is what a gallery shows and what a limit
		// should keep.
		order, err := e.String("date_added DESC")
		if err != nil {
			return err
		}
		cursor, collection, err := queryCursor(e, kind, projection, where, args, order)
		if err != nil || cursor.IsNil() {
			return err
		}
		defer e.InvokeVoid(cursor, "close", jni.Sig(jni.TVoid))

		out, err = readItems(e, cursor, collection, modern, q.Limit)
		return err
	})
	return out, err
}
