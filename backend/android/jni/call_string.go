//go:build android

package jni

// CallString is CallObject and then reading the String, which is common
// enough to be worth one call instead of three lines.
func (e *Env) CallString(o Object, m Method, args ...Value) (string, error) {
	r, err := e.CallObject(o, m, args...)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}
