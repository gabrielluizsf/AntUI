//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func mustString(e *jni.Env, s string) jni.Value {
	o, err := e.String(s)
	if err != nil {
		return jni.Ref(jni.Object{})
	}
	return jni.Ref(o)
}
