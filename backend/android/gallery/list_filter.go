//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// listFilter builds the selection and its arguments: the app's own media, a
// folder, or nothing at all.
func listFilter(e *jni.Env, q Query, modern bool) (where, args jni.Object, err error) {
	switch {
	case q.Mine:
		// Everything MediaStore reports to an app with no permission is
		// already its own, but saying so makes the intent plain and keeps
		// the answer the same when a permission is later added.
		if where, err = e.String("owner_package_name = ?"); err != nil {
			return jni.Object{}, jni.Object{}, err
		}
		name, err := e.InvokeString(app.Context(), "getPackageName", jni.Sig(jni.TString))
		if err != nil {
			return jni.Object{}, jni.Object{}, err
		}
		if args, err = e.NewStrings([]string{name}); err != nil {
			return jni.Object{}, jni.Object{}, err
		}
	case q.Album != "" && modern:
		if where, err = e.String(colPath + " LIKE ?"); err != nil {
			return jni.Object{}, jni.Object{}, err
		}
		if args, err = e.NewStrings([]string{"%/" + q.Album + "/%"}); err != nil {
			return jni.Object{}, jni.Object{}, err
		}
	}
	return where, args, nil
}
