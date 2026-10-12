//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func geocoder(e *jni.Env) (jni.Object, error) {
	return e.Make("android/location/Geocoder",
		jni.Sig(jni.TVoid, jni.TContext), jni.Ref(app.Context()))
}
