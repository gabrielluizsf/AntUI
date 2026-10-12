//go:build android

package torch

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// cameraWithFlash finds the first camera with a flash, once. The list does
// not change while an app runs, and asking costs a walk over every camera's
// characteristics.
func cameraWithFlash() (string, error) {
	idOnce.Do(func() {
		idErr = jni.Do(func(e *jni.Env) error {
			id, err := findFlash(e)
			if err != nil {
				return err
			}
			flashID = id
			return nil
		})
	})
	if idErr != nil {
		return "", idErr
	}
	return flashID, nil
}
