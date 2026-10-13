//go:build android

package prefs

import "github.com/gabrielluizsf/antui/backend/android/jni"

// SetBool writes a true-or-false setting.
func (p *Prefs) SetBool(key string, value bool) error {
	return p.putScalar(key, "putBoolean", jni.TBool, jni.Bool(value))
}
