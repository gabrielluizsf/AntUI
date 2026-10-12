//go:build android

package dialog

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// dialogButtons adds the buttons that have a label, each answering through
// the one listener.
func dialogButtons(e *jni.Env, builder, listener jni.Object, positive, negative, neutral string) error {
	for _, b := range []struct{ method, label string }{
		{"setPositiveButton", positive},
		{"setNegativeButton", negative},
		{"setNeutralButton", neutral},
	} {
		if b.label == "" {
			continue
		}
		js, err := e.String(b.label)
		if err != nil {
			return err
		}
		if _, err := e.Invoke(builder, b.method,
			jni.Sig(builderType, jni.TClass("java/lang/CharSequence"),
				jni.TClass("android/content/DialogInterface$OnClickListener")),
			jni.Ref(js), jni.Ref(listener)); err != nil {
			return err
		}
	}
	return nil
}
