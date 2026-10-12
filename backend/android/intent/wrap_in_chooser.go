//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func wrapInChooser(e *jni.Env, in jni.Object, title string) (jni.Object, error) {
	jtitle, err := e.String(title)
	if err != nil {
		return jni.Object{}, err
	}
	return e.Static("android/content/Intent", "createChooser",
		jni.Sig(jni.TClass("android/content/Intent"),
			jni.TClass("android/content/Intent"), jni.TClass("java/lang/CharSequence")),
		jni.Ref(in), jni.Ref(jtitle))
}
