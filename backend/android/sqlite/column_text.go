//go:build android

package sqlite

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func (c column) text() (string, error) {
	return c.e.InvokeString(c.cursor, "getString", jni.Sig(jni.TString, jni.TInt), jni.Int(c.i))
}
