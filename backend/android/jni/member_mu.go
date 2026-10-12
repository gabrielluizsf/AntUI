//go:build android

package jni

import (
	"sync"
)

var (
	memberMu sync.RWMutex
	methods  = map[memberKey]Method{}
	fields   = map[memberKey]Field{}
)
