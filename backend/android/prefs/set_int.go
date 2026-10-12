//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// SetInt writes an int setting.
func (p *Prefs) SetInt(key string, value int32) error {
	return p.putScalar(key, "putInt", jni.TInt, jni.Int(value))
}
