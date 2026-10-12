//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// setIcon: the icon is a resource in the package, chosen when the APK is
// built and not at runtime.
func (b *androidWindow) setIcon([]*canvas.Canvas) bool { return false }
