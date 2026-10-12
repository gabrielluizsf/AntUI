//go:build android

package biometric

import (
	"errors"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// plainly turns the one failure everybody hits into a sentence.
func plainly(err error) error {
	var je *jni.Error
	if errors.As(err, &je) &&
		je.Class == "java.lang.SecurityException" &&
		strings.Contains(je.Message, "USE_BIOMETRIC") {
		return ErrNotDeclared
	}
	return err
}
