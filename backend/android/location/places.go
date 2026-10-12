//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// places reads a java.util.List of Address into Go.
func places(e *jni.Env, list jni.Object) ([]Place, error) {
	if list.IsNil() {
		return nil, nil
	}
	n, err := e.InvokeInt(list, "size", jni.Sig(jni.TInt))
	if err != nil || n == 0 {
		return nil, err
	}
	out := make([]Place, 0, n)
	err = e.Frame(int(n)*8+16, func() error {
		for i := range int(n) {
			item, err := e.Invoke(list, "get", jni.Sig(jni.TObject, jni.TInt), jni.Int(int32(i)))
			if err != nil || item.IsNil() {
				continue
			}
			p, err := place(e, item)
			if err != nil {
				return err
			}
			out = append(out, p)
		}
		return nil
	})
	return out, err
}
