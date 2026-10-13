//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// Assets is the reader for the files packed into the APK.
func (a *App) Assets() *ndk.AssetManager {
	if a == nil {
		return nil
	}
	return ndk.WrapAssets(a.assets)
}
