//go:build android

package jni

// IsNil reports whether the class was never found.
func (c Class) IsNil() bool { return c.ref == 0 }
