//go:build android

package assets

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

func manager() (*ndk.AssetManager, error) {
	a := app.Current()
	if a == nil {
		return nil, app.ErrNoUIThread
	}
	m := a.Assets()
	if m == nil {
		return nil, errors.New("assets: this app has no asset manager")
	}
	return m, nil
}
