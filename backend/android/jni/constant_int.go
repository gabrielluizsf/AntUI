//go:build android

package jni

// ConstantInt reads an int constant off a class — the API levels,
// the flags and the modes Android keeps as static final fields.
func (e *Env) ConstantInt(class, name string) (int32, error) {
	c, err := e.Class(class)
	if err != nil {
		return 0, err
	}
	f, err := e.cachedField(c, name, TInt, kindStaticField)
	if err != nil {
		return 0, err
	}
	return e.StaticInt(c, f)
}
