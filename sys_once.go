package antui

import (
	"sync"
)

var (
	sysOnce sync.Once
	sysData []byte
	sysPath string
	sysErr  error
)
