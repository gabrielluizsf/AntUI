//go:build android

package assets

import (
	"time"
)

func (i info) ModTime() time.Time { return time.Time{} }
