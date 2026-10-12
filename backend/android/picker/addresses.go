//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// addresses reads the one address or the list of them, because the picker
// returns them in two different places depending on how many were chosen —
// even when multiple were allowed and the user picked one.
func addresses(data jni.Object) ([]string, error) {
	var out []string
	err := jni.Do(func(e *jni.Env) error {
		clip, err := e.Invoke(data, "getClipData",
			jni.Sig(jni.TClass("android/content/ClipData")))
		if err != nil {
			return err
		}
		if !clip.IsNil() {
			n, err := e.InvokeInt(clip, "getItemCount", jni.Sig(jni.TInt))
			if err != nil {
				return err
			}
			return e.Frame(int(n)*4+16, func() error {
				for i := range int(n) {
					item, err := e.Invoke(clip, "getItemAt",
						jni.Sig(jni.TClass("android/content/ClipData$Item"), jni.TInt),
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
		}
		uri, err := e.Invoke(data, "getData", jni.Sig(jni.TClass("android/net/Uri")))
		if err != nil || uri.IsNil() {
			return err
		}
		s, err := e.InvokeString(uri, "toString", jni.Sig(jni.TString))
		if err == nil && s != "" {
			out = append(out, s)
		}
		return err
	})
	return out, err
}
