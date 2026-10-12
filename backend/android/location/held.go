//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

// held reports whether either location permission has been granted.
func held() (bool, error) {
	for _, p := range []string{permission.FineLocation, permission.CoarseLocation} {
		ok, err := permission.Held(p)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
