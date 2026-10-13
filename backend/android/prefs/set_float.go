//go:build android

package prefs

import "github.com/gabrielluizsf/antui/backend/android/jni"

// SetFloat writes a float setting.
func (p *Prefs) SetFloat(key string, value float32) error {
	return p.putScalar(key, "putFloat", jni.TFloat, jni.Float(value))
}
