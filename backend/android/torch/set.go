//go:build android

package torch

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Set lights the lamp or puts it out.
func Set(on bool) error {
	id, err := cameraWithFlash()
	if err != nil {
		return err
	}
	return jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, "camera")
		if err != nil {
			return err
		}
		jid, err := e.String(id)
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "setTorchMode",
			jni.Sig(jni.TVoid, jni.TString, jni.TBool), jni.Ref(jid), jni.Bool(on))
	})
}
