//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// builderText sets the title and the body.
func builderText(e *jni.Env, builder jni.Object, n Notification) error {
	title, err := e.String(n.Title)
	if err != nil {
		return err
	}
	text, err := e.String(n.Text)
	if err != nil {
		return err
	}
	cs := jni.TClass("java/lang/CharSequence")
	b := jni.TClass("android/app/Notification$Builder")
	if err := setBuilder(e, builder, "setContentTitle", jni.Sig(b, cs), jni.Ref(title)); err != nil {
		return err
	}
	return setBuilder(e, builder, "setContentText", jni.Sig(b, cs), jni.Ref(text))
}
