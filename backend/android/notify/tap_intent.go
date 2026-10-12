//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// tapIntent builds what happens when the notification is tapped: the app is
// brought to the front, and the intent reaches it through onNewIntent.
func tapIntent(e *jni.Env, id int32, extra map[string]string) (jni.Object, error) {
	cls, err := e.Class(app.ShimClass)
	if err != nil {
		// Built without the shim: there is nothing that can receive the tap,
		// so the notification is posted without one rather than not at all.
		return jni.Object{}, nil
	}
	intent, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TContext, jni.TClass("java/lang/Class")),
		jni.Ref(app.Context()), jni.Ref(cls.Object()))
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(intent, "setFlags",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
		jni.Int(flagSingleTop|flagClearTop)); err != nil {
		return jni.Object{}, err
	}
	for k, v := range extra {
		jk, err := e.String(k)
		if err != nil {
			return jni.Object{}, err
		}
		jv, err := e.String(v)
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(intent, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TString),
			jni.Ref(jk), jni.Ref(jv)); err != nil {
			return jni.Object{}, err
		}
	}
	return e.Static("android/app/PendingIntent", "getActivity",
		jni.Sig(jni.TClass("android/app/PendingIntent"),
			jni.TContext, jni.TInt, jni.TClass("android/content/Intent"), jni.TInt),
		jni.Ref(app.Context()), jni.Int(id), jni.Ref(intent),
		jni.Int(flagUpdateCurrent|flagImmutable))
}
