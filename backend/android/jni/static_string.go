//go:build android

package jni

// StaticString reads a String field of a class, as Go text.
func (e *Env) StaticString(c Class, f Field) (string, error) {
	r, err := e.StaticObject(c, f)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}
