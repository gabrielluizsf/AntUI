//go:build android

package jni

// ConstantString and ConstantInt are Constant for the two types most of them
// are.
func (e *Env) ConstantString(class, name string) (string, error) {
	o, err := e.Constant(class, name, TString)
	if err != nil {
		return "", err
	}
	return e.GoString(o)
}
