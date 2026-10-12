//go:build android

package intent

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// start hands the intent to the system, and turns the one failure everybody
// hits into a sentence rather than a Java class name.
func start(e *jni.Env, in jni.Object) error {
	err := e.InvokeVoid(app.Context(), "startActivity",
		jni.Sig(jni.TVoid, jni.TClass("android/content/Intent")), jni.Ref(in))
	var je *jni.Error
	if errorsAs(err, &je) && je.Class == "android.content.ActivityNotFoundException" {
		return fmt.Errorf("intent: nothing on this device can open that")
	}
	return err
}
