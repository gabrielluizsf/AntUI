//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Geocode turns an address into places. max bounds how many come back; one
// is usually right, and more is for offering a choice.
//
// **It goes to the network and blocks.** Not from the frame loop, and not
// from anything holding a lock a frame wants.
func Geocode(query string, max int) ([]Place, error) {
	if max <= 0 {
		max = 1
	}
	var out []Place
	err := jni.Do(func(e *jni.Env) error {
		coder, err := geocoder(e)
		if err != nil {
			return err
		}
		q, err := e.String(query)
		if err != nil {
			return err
		}
		list, err := e.Invoke(coder, "getFromLocationName",
			jni.Sig(jni.TClass("java/util/List"), jni.TString, jni.TInt),
			jni.Ref(q), jni.Int(int32(max)))
		if err != nil {
			return err
		}
		out, err = places(e, list)
		return err
	})
	return out, err
}
