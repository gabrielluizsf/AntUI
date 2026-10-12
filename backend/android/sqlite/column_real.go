//go:build android

package sqlite

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func (c column) real() (float64, error) {
	v, err := c.e.InvokeFloat(c.cursor, "getFloat", jni.Sig(jni.TFloat, jni.TInt), jni.Int(c.i))
	return float64(v), err
}
