//go:build android

package jni

// CallStaticString calls a class method that returns a String, and
// reads it out as Go text.
func (e *Env) CallStaticString(c Class, m Method, args ...Value) (string, error) {
	r, err := e.CallStaticObject(c, m, args...)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}
