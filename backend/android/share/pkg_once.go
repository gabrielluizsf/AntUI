//go:build android

package share

import (
	"sync"
)

var (
	pkgOnce sync.Once
	pkgName string
)
