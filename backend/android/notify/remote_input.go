//go:build android

package notify

import "github.com/gabrielluizsf/antui/backend/android/jni"

func remoteInput(e *jni.Env, key, label string) (jni.Object, error) {
	jkey, err := e.String(key)
	if err != nil {
		return jni.Object{}, err
	}
	b, err := e.Make("android/app/RemoteInput$Builder",
		jni.Sig(jni.TVoid, jni.TString), jni.Ref(jkey))
	if err != nil {
		return jni.Object{}, err
	}
	if label != "" {
		jlabel, err := e.String(label)
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(b, "setLabel",
			jni.Sig(jni.TClass("android/app/RemoteInput$Builder"),
				jni.TClass("java/lang/CharSequence")), jni.Ref(jlabel)); err != nil {
			return jni.Object{}, err
		}
	}
	return e.Invoke(b, "build", jni.Sig(jni.TClass("android/app/RemoteInput")))
}
