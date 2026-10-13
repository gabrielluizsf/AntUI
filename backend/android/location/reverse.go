//go:build android

package location

import "github.com/gabrielluizsf/antui/backend/android/jni"

// Reverse turns a position into addresses.
//
// The same warning: it goes to the network and blocks.
func Reverse(latitude, longitude float64, max int) ([]Place, error) {
	if max <= 0 {
		max = 1
	}
	var out []Place
	err := jni.Do(func(e *jni.Env) error {
		coder, err := geocoder(e)
		if err != nil {
			return err
		}
		list, err := e.Invoke(coder, "getFromLocation",
			jni.Sig(jni.TClass("java/util/List"), jni.TDouble, jni.TDouble, jni.TInt),
			jni.Double(latitude), jni.Double(longitude), jni.Int(int32(max)))
		if err != nil {
			return err
		}
		out, err = places(e, list)
		return err
	})
	return out, err
}
