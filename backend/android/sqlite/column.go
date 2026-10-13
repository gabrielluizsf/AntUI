//go:build android

package sqlite

import "github.com/gabrielluizsf/antui/backend/android/jni"

// column is one value of one row, with the cursor it is read from.
type column struct {
	e      *jni.Env
	cursor jni.Object
	i      int32
}
