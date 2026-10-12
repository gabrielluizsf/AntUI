//go:build android

package biometric

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// promptError turns the platform's numbered reason into one of ours.
func promptError(e *jni.Env, args []jni.Object) error {
	if len(args) == 0 {
		return ErrCancelled
	}
	code, err := e.InvokeInt(args[0], "intValue", jni.Sig(jni.TInt))
	if err != nil {
		return ErrCancelled
	}
	switch code {
	case 5, 10: // ERROR_CANCELED, ERROR_USER_CANCELED
		return ErrCancelled
	case 13: // ERROR_NEGATIVE_BUTTON
		return ErrRefused
	case 7, 9: // ERROR_LOCKOUT, ERROR_LOCKOUT_PERMANENT
		return ErrLockedOut
	}
	message := ""
	if len(args) > 1 {
		message, _ = e.GoString(args[1])
	}
	if message == "" {
		message = fmt.Sprintf("error %d", code)
	}
	return fmt.Errorf("biometric: %s", message)
}
