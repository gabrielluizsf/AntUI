//go:build android

package jni

// GoStrings reads a String[] into Go.
func (e *Env) GoStrings(o Object) ([]string, error) {
	n, err := e.Len(o)
	if err != nil || n == 0 {
		return nil, err
	}
	out := make([]string, n)
	err = e.Frame(n+8, func() error {
		for i := range n {
			el, err := e.Index(o, i)
			if err != nil {
				return err
			}
			if out[i], err = e.GoString(el); err != nil {
				return err
			}
		}
		return nil
	})
	return out, err
}
