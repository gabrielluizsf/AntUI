//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/jni"

// Listener makes one of the shim's own listener classes — the ones that had
// to be written out because Android wanted an abstract class rather than an
// interface.
func Listener(e *jni.Env, class string, token int64) (jni.Object, error) {
	return e.Static(class, "create", jni.Sig(jni.TObject, jni.TLong), jni.Long(token))
}
