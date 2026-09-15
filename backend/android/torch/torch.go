//go:build android

// Package torch turns the camera's flash on and steady, as a lamp.
//
// It needs **no permission**, which is worth saying because everything else
// to do with the camera does. Since API 23 the platform treats the torch as
// its own thing: an app may light it without being allowed to see through
// the lens.
package torch

import (
	"errors"
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// ErrNoFlash is a device with no lamp on any of its cameras.
var ErrNoFlash = errors.New("torch: no camera on this device has a flash")

var (
	idOnce  sync.Once
	flashID string
	idErr   error
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

// Available reports whether this device has a lamp to light.
func Available() bool {
	_, err := cameraWithFlash()
	return err == nil
}

// cameraWithFlash finds the first camera with a flash, once. The list does
// not change while an app runs, and asking costs a walk over every camera's
// characteristics.
func cameraWithFlash() (string, error) {
	idOnce.Do(func() {
		idErr = jni.Do(func(e *jni.Env) error {
			manager, err := app.Service(e, "camera")
			if err != nil {
				return err
			}
			ids, err := e.Invoke(manager, "getCameraIdList",
				jni.Sig(jni.TArray(jni.TString)))
			if err != nil {
				return err
			}
			names, err := e.GoStrings(ids)
			if err != nil {
				return err
			}
			// FLASH_INFO_AVAILABLE is a key object, not a name: the
			// characteristics are a typed map and the keys are constants on
			// the class.
			key, err := e.Constant("android/hardware/camera2/CameraCharacteristics",
				"FLASH_INFO_AVAILABLE",
				jni.TClass("android/hardware/camera2/CameraCharacteristics$Key"))
			if err != nil {
				return err
			}
			for _, id := range names {
				jid, err := e.String(id)
				if err != nil {
					return err
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
					flashID = id
					return nil
				}
			}
			return ErrNoFlash
		})
	})
	if idErr != nil {
		return "", idErr
	}
	return flashID, nil
}
