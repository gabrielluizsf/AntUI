//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// builderPost builds the notification and hands it to the manager.
func builderPost(e *jni.Env, builder jni.Object, n Notification) error {
	built, err := e.Invoke(builder, "build",
		jni.Sig(jni.TClass("android/app/Notification")))
	if err != nil {
		return err
	}
	manager, err := app.Service(e, app.ServiceNotification)
	if err != nil {
		return err
	}
	return e.InvokeVoid(manager, "notify",
		jni.Sig(jni.TVoid, jni.TInt, jni.TClass("android/app/Notification")),
		jni.Int(n.ID), jni.Ref(built))
}
