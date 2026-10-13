//go:build android

package sqlite

import "github.com/gabrielluizsf/antui/backend/android/jni"

// Rows is the answer to a query, walked one row at a time.
type Rows struct {
	cursor  jni.Object
	columns []string
	err     error
	closed  bool
}
