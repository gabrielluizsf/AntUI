//go:build android

package speech

import "github.com/gabrielluizsf/antui/backend/android/jni"

func startListening(e *jni.Env, recogniser, intent jni.Object) error {
	return e.InvokeVoid(recogniser, "startListening",
		jni.Sig(jni.TVoid, jni.TClass("android/content/Intent")), jni.Ref(intent))
}
