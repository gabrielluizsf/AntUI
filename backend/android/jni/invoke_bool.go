//go:build android

package jni

// InvokeBool is [Env.Invoke] for a method that returns a boolean.
func (e *Env) InvokeBool(o Object, name, sig string, args ...Value) (bool, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return false, err
	}
	return e.CallBool(o, m, args...)
}
