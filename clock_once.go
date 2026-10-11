package antui

import (
	"sync"
	"time"
)

var (
	clockOnce sync.Once
	clockZero time.Time
)
