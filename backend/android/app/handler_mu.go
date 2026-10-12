//go:build android

package app

import (
	"sync"
	"sync/atomic"
)

// The registry. Java holds a number and nothing else: a Java field holding a
// Go pointer is something neither side's collector could keep honest.
var (
	handlerMu sync.Mutex
	handlers  = map[int64]Handler{}
	lastToken atomic.Int64
)
