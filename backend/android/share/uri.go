//go:build android

package share

import (
	"net/url"

	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// URI is the address of a file already in [Dir].
//
// It is a content:// address and not a path. **Since Android 7 a path is not
// something one app may give another**: an intent carrying a file:// address
// throws, deliberately, because a path is meaningless in another app's world
// and used to be a way to hand out things that were never meant to leave.
func URI(name string) (string, error) {
	pkg, err := packageName()
	if err != nil {
		return "", err
	}
	return "content://" + pkg + shim.ProviderSuffix + "/" + url.PathEscape(name), nil
}
