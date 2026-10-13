//go:build android

package dialog

import "github.com/gabrielluizsf/antui/backend/android/jni"

// dialogText sets one text property on the builder, skipping an empty one.
func dialogText(e *jni.Env, builder jni.Object, method, value string) error {
	if value == "" {
		return nil
	}
	js, err := e.String(value)
	if err != nil {
		return err
	}
	_, err = e.Invoke(builder, method,
		jni.Sig(builderType, jni.TClass("java/lang/CharSequence")), jni.Ref(js))
	return err
}
