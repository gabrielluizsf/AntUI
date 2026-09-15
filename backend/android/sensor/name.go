//go:build android

package sensor

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

var (
	nameOnce sync.Once
	pkgName  string
)

// packageName is what the sensor manager is opened for. It is asked once:
// it never changes, and asking costs a trip into Java.
func packageName() string {
	nameOnce.Do(func() {
		jni.Do(func(e *jni.Env) error {
			pkgName, _ = e.InvokeString(app.Context(), "getPackageName", jni.Sig(jni.TString))
			return nil
		})
	})
	return pkgName
}
