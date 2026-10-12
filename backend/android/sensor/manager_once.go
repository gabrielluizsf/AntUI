//go:build android

package sensor

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

var (
	managerOnce sync.Once
	manager     *ndk.SensorManager
	managerErr  error
)
