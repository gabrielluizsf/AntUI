//go:build android

package jni

// Invoke calls a method on an object, finding the method by name.
func (e *Env) Invoke(o Object, name, sig string, args ...Value) (Object, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return Object{}, err
	}
	return e.CallObject(o, m, args...)
}
