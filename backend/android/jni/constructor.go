//go:build android

package jni

// Constructor finds a constructor. Its name is "<init>" and it returns void,
// which is a rule rather than a convention, so the signature only needs its
// arguments: Constructor(cls, Sig(TVoid, TInt)).
func (e *Env) Constructor(c Class, sig string) (Method, error) {
	return e.Method(c, "<init>", sig)
}
