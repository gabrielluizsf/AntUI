//go:build android

package notify

import (
	"sync"
)

var (
	actionMu   sync.Mutex
	actionFns  = map[string]func(id int32, reply string){}
	actionTok  int64
	actionOnce sync.Once
)
