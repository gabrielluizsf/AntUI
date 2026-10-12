//go:build android

package jni

// NewStrings makes a String[] out of Go strings.
func (e *Env) NewStrings(v []string) (Object, error) {
	cls, err := e.Class("java/lang/String")
	if err != nil {
		return Object{}, err
	}
	arr, err := e.NewObjects(cls, len(v))
	if err != nil {
		return Object{}, err
	}
	// A frame of its own: each string is a local reference, and a list of a
	// few hundred would otherwise fill the table the caller is sharing.
	err = e.Frame(len(v)+8, func() error {
		for i, s := range v {
			js, err := e.String(s)
			if err != nil {
				return err
			}
			if err := e.SetIndex(arr, i, js); err != nil {
				return err
			}
		}
		return nil
	})
	return arr, err
}
