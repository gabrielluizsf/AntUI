//go:build android

package jni

// InvokeLong is [Env.Invoke] for a method that returns a long.
func (e *Env) InvokeLong(o Object, name, sig string, args ...Value) (int64, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return 0, err
	}
	return e.CallLong(o, m, args...)
}
