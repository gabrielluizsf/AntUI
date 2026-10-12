//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Kept is every address the app is still allowed to reach, from choices the
// user made in earlier runs.
func Kept() ([]string, error) {
	var out []string
	err := jni.Do(func(e *jni.Env) error {
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		list, err := e.Invoke(resolver, "getPersistedUriPermissions",
			jni.Sig(jni.TClass("java/util/List")))
		if err != nil || list.IsNil() {
			return err
		}
		n, err := e.InvokeInt(list, "size", jni.Sig(jni.TInt))
		if err != nil || n == 0 {
			return err
		}
		return e.Frame(int(n)*4+16, func() error {
			for i := range int(n) {
				item, err := e.Invoke(list, "get", jni.Sig(jni.TObject, jni.TInt),
					jni.Int(int32(i)))
				if err != nil || item.IsNil() {
					continue
				}
				uri, err := e.Invoke(item, "getUri",
					jni.Sig(jni.TClass("android/net/Uri")))
				if err != nil || uri.IsNil() {
					continue
				}
				s, err := e.InvokeString(uri, "toString", jni.Sig(jni.TString))
				if err == nil && s != "" {
					out = append(out, s)
				}
			}
			return nil
		})
	})
	return out, err
}
