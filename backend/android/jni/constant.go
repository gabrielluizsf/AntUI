//go:build android

package jni

// Constant reads a static field, which is where the platform keeps most of
// its named values — Context.VIBRATOR_SERVICE, Intent.ACTION_SEND, and a
// thousand others.
func (e *Env) Constant(class, name, sig string) (Object, error) {
	c, err := e.Class(class)
	if err != nil {
		return Object{}, err
	}
	f, err := e.cachedField(c, name, sig, kindStaticField)
	if err != nil {
		return Object{}, err
	}
	return e.StaticObject(c, f)
}
