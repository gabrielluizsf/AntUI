//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// imeInset reads the keyboard's height on the platforms that report it
// separately. It answers 0 rather than an error on the ones that do not,
// because "there is no such call here" is not a failure.
func imeInset(e *jni.Env, ins jni.Object, wiCls jni.Class) int {
	typeCls, err := e.Class("android/view/WindowInsets$Type")
	if err != nil {
		return 0
	}
	imeM, err := e.StaticMethod(typeCls, "ime", jni.Sig(jni.TInt))
	if err != nil {
		return 0
	}
	which, err := e.CallStaticInt(typeCls, imeM)
	if err != nil {
		return 0
	}
	getInsets, err := e.Method(wiCls, "getInsets",
		jni.Sig(jni.TClass("android/graphics/Insets"), jni.TInt))
	if err != nil {
		return 0
	}
	got, err := e.CallObject(ins, getInsets, jni.Int(which))
	if err != nil || got.IsNil() {
		return 0
	}
	insetsCls, err := e.Class("android/graphics/Insets")
	if err != nil {
		return 0
	}
	f, err := e.Field(insetsCls, "bottom", jni.TInt)
	if err != nil {
		return 0
	}
	v, err := e.GetInt(got, f)
	if err != nil {
		return 0
	}
	return int(v)
}
