//go:build android

package display

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// All is every screen the device can see, which is more than one while it is
// casting or plugged into a monitor.
func All() ([]Display, error) {
	var out []Display
	err := jni.Do(func(e *jni.Env) error {
		dm, err := app.Service(e, app.ServiceDisplay)
		if err != nil {
			return err
		}
		arr, err := e.Invoke(dm, "getDisplays",
			jni.Sig(jni.TArray(jni.TClass("android/view/Display"))))
		if err != nil || arr.IsNil() {
			return err
		}
		n, err := e.Len(arr)
		if err != nil {
			return err
		}
		out = make([]Display, 0, n)
		return e.Frame(n+8, func() error {
			for i := range n {
				d, err := e.Index(arr, i)
				if err != nil {
					return err
				}
				one, err := read(e, d)
				if err != nil {
					return err
				}
				out = append(out, one)
			}
			return nil
		})
	})
	return out, err
}
