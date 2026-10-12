//go:build android

package jni

// InvokeInt is [Env.Invoke] for a method that returns an int.
func (e *Env) InvokeInt(o Object, name, sig string, args ...Value) (int32, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return 0, err
	}
	return e.CallInt(o, m, args...)
}
