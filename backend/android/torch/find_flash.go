//go:build android

package torch

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// findFlash walks the cameras and returns the first id whose flash is
// available.
func findFlash(e *jni.Env) (string, error) {
	manager, err := app.Service(e, "camera")
	if err != nil {
		return "", err
	}
	ids, err := e.Invoke(manager, "getCameraIdList",
		jni.Sig(jni.TArray(jni.TString)))
	if err != nil {
		return "", err
	}
	names, err := e.GoStrings(ids)
	if err != nil {
		return "", err
	}
	// FLASH_INFO_AVAILABLE is a key object, not a name: the
	// characteristics are a typed map and the keys are constants on
	// the class.
	key, err := e.Constant("android/hardware/camera2/CameraCharacteristics",
		"FLASH_INFO_AVAILABLE",
		jni.TClass("android/hardware/camera2/CameraCharacteristics$Key"))
	if err != nil {
		return "", err
	}
	for _, id := range names {
		jid, err := e.String(id)
		if err != nil {
			return "", err
		}
		chars, err := e.Invoke(manager, "getCameraCharacteristics",
			jni.Sig(jni.TClass("android/hardware/camera2/CameraCharacteristics"),
				jni.TString), jni.Ref(jid))
		if err != nil || chars.IsNil() {
			continue
		}
		has, err := e.Invoke(chars, "get",
			jni.Sig(jni.TObject,
				jni.TClass("android/hardware/camera2/CameraCharacteristics$Key")),
			jni.Ref(key))
		if err != nil || has.IsNil() {
			continue
		}
		on, err := e.InvokeBool(has, "booleanValue", jni.Sig(jni.TBool))
		if err == nil && on {
			return id, nil
		}
	}
	return "", ErrNoFlash
}
