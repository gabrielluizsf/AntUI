//go:build android

package sensor

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

func sensors() (*ndk.SensorManager, error) {
	managerOnce.Do(func() {
		name := ""
		if a := app.Current(); a != nil {
			// The package name, which the platform uses from API 26 to
			// decide which sensors this app may see.
			name = packageName()
		}
		manager, managerErr = ndk.Sensors(name)
	})
	return manager, managerErr
}
