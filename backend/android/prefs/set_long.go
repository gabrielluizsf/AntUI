//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// SetLong writes a 64-bit setting.
func (p *Prefs) SetLong(key string, value int64) error {
	return p.putScalar(key, "putLong", jni.TLong, jni.Long(value))
}
