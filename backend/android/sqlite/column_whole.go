//go:build android

package sqlite

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func (c column) whole() (int64, error) {
	return c.e.InvokeLong(c.cursor, "getLong", jni.Sig(jni.TLong, jni.TInt), jni.Int(c.i))
}
