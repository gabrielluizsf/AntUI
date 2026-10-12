//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func errorsAs(err error, target **jni.Error) bool {
	je, ok := err.(*jni.Error)
	if ok {
		*target = je
	}
	return ok
}
