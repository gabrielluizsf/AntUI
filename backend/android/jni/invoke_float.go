//go:build android

package jni

// InvokeFloat is [Env.Invoke] for a method that returns a float.
func (e *Env) InvokeFloat(o Object, name, sig string, args ...Value) (float32, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return 0, err
	}
	return e.CallFloat(o, m, args...)
}
