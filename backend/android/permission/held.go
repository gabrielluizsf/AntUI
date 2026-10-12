//go:build android

package permission

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
)

// Held reports whether a permission has already been granted.
func Held(name string) (bool, error) { return app.HasPermission(name) }
