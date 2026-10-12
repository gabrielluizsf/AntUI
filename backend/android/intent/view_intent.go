//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func viewIntent(e *jni.Env, address, mime string) (jni.Object, error) {
	action, err := e.String(ActionView)
	if err != nil {
		return jni.Object{}, err
	}
	jaddr, err := e.String(address)
	if err != nil {
		return jni.Object{}, err
	}
	uri, err := e.Static("android/net/Uri", "parse",
		jni.Sig(jni.TClass("android/net/Uri"), jni.TString), jni.Ref(jaddr))
	if err != nil {
		return jni.Object{}, err
	}
	in, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TString, jni.TClass("android/net/Uri")),
		jni.Ref(action), jni.Ref(uri))
	if err != nil {
		return jni.Object{}, err
	}
	if mime != "" {
		jmime, err := e.String(mime)
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(in, "setDataAndType",
			jni.Sig(jni.TClass("android/content/Intent"),
				jni.TClass("android/net/Uri"), jni.TString),
			jni.Ref(uri), jni.Ref(jmime)); err != nil {
			return jni.Object{}, err
		}
	}
	return in, nil
}
